package raft

// Implementacao da parte 2A: eleicao e heartbeats. A API do laboratorio
// e preservada; replicacao (2B) e persistencia (2C) estao fora do escopo.
import (
	"labrpc"
	"math/rand"
	"sync"
	"time"
)

type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool
	Snapshot    []byte
}

type LogEntry struct {
	Term    int
	Command interface{}
}

type role int

const (
	follower role = iota
	candidate
	leader
)

const (
	heartbeatInterval = 120 * time.Millisecond
	electionMinimum   = 450 * time.Millisecond
	electionJitter    = 450 * time.Millisecond
	pollInterval      = 10 * time.Millisecond
)

type Raft struct {
	mu        sync.Mutex
	peers     []*labrpc.ClientEnd
	persister *Persister
	me        int

	currentTerm int
	votedFor    int
	state       role
	log         []LogEntry // indice zero e uma sentinela; vazio na parte 2A
	votes       int

	rng               *rand.Rand // acesso somente sob mu
	electionDeadline  time.Time
	heartbeatDeadline time.Time
	stopped           bool
	stopCh            chan struct{}
}

func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.currentTerm, !rf.stopped && rf.state == leader
}

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

// Todos os auxiliares com sufixo Locked exigem rf.mu adquirido.
func (rf *Raft) resetElectionLocked() {
	delay := electionMinimum + time.Duration(rf.rng.Int63n(int64(electionJitter)))
	rf.electionDeadline = time.Now().Add(delay)
}

// Um termo maior invalida o voto anterior. Nao reiniciamos o timer aqui:
// receber um pedido de voto recusado nao deve adiar uma eleicao.
func (rf *Raft) observeTermLocked(term int) {
	if term > rf.currentTerm {
		rf.currentTerm = term
		rf.votedFor = -1
		rf.state = follower
	}
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.VoteGranted = false
	if !rf.stopped {
		rf.observeTermLocked(args.Term)
	}
	reply.Term = rf.currentTerm
	if rf.stopped || args.Term < rf.currentTerm {
		return
	}
	last := len(rf.log) - 1
	upToDate := args.LastLogTerm > rf.log[last].Term ||
		(args.LastLogTerm == rf.log[last].Term && args.LastLogIndex >= last)
	if upToDate && (rf.votedFor == -1 || rf.votedFor == args.CandidateId) {
		rf.votedFor = args.CandidateId
		reply.VoteGranted = true
		rf.resetElectionLocked()
	}
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.Success = false
	if !rf.stopped {
		rf.observeTermLocked(args.Term)
	}
	reply.Term = rf.currentTerm
	if rf.stopped || args.Term < rf.currentTerm {
		return
	}
	// Inclui candidato que recebe heartbeat do mesmo termo. O voto desse
	// termo e preservado, evitando conceder dois votos a candidatos distintos.
	rf.state = follower
	rf.resetElectionLocked()
	if args.PrevLogIndex < 0 || args.PrevLogIndex >= len(rf.log) ||
		rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		return
	}
	// Apenas heartbeats neste trabalho; nao simular replicacao bem-sucedida.
	if len(args.Entries) != 0 {
		return
	}
	reply.Success = true
}

func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	return rf.peers[server].Call("Raft.RequestVote", args, reply)
}

func (rf *Raft) becomeLeaderLocked() {
	rf.state = leader
	rf.broadcastHeartbeatLocked() // primeiro heartbeat imediatamente apos maioria
}

func (rf *Raft) startElectionLocked() {
	rf.state = candidate
	rf.currentTerm++
	rf.votedFor = rf.me
	rf.votes = 1
	rf.resetElectionLocked()
	last := len(rf.log) - 1
	args := RequestVoteArgs{rf.currentTerm, rf.me, last, rf.log[last].Term}
	if rf.votes > len(rf.peers)/2 {
		rf.becomeLeaderLocked() // cluster de um servidor
		return
	}
	for peer := range rf.peers {
		if peer == rf.me {
			continue
		}
		go func(server int, request RequestVoteArgs) {
			var reply RequestVoteReply
			if !rf.sendRequestVote(server, &request, &reply) {
				return
			}
			rf.mu.Lock()
			defer rf.mu.Unlock()
			if rf.stopped {
				return
			}
			// Mesmo resposta de uma eleicao antiga pode revelar termo maior.
			rf.observeTermLocked(reply.Term)
			if rf.state != candidate || rf.currentTerm != request.Term ||
				reply.Term != request.Term || !reply.VoteGranted {
				return
			}
			// Uma unica RPC por destinatario nesta rodada: cada voto conta uma vez.
			rf.votes++
			if rf.votes > len(rf.peers)/2 {
				rf.becomeLeaderLocked()
			}
		}(peer, args)
	}
}

func (rf *Raft) broadcastHeartbeatLocked() {
	rf.heartbeatDeadline = time.Now().Add(heartbeatInterval)
	// O log permanece vazio em 2A; o prefixo consistente e a sentinela 0/0.
	args := AppendEntriesArgs{Term: rf.currentTerm, LeaderId: rf.me}
	for peer := range rf.peers {
		if peer == rf.me {
			continue
		}
		go func(server int, request AppendEntriesArgs) {
			var reply AppendEntriesReply
			// A chamada ocorre fora de mu: uma rede lenta nao bloqueia o peer.
			if !rf.peers[server].Call("Raft.AppendEntries", &request, &reply) {
				return
			}
			rf.mu.Lock()
			defer rf.mu.Unlock()
			if !rf.stopped {
				rf.observeTermLocked(reply.Term)
			}
		}(peer, args)
	}
}

func (rf *Raft) ticker() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-rf.stopCh:
			return
		case <-ticker.C:
			rf.mu.Lock()
			if rf.stopped {
				rf.mu.Unlock()
				return
			}
			now := time.Now()
			if rf.state == leader {
				if !now.Before(rf.heartbeatDeadline) {
					rf.broadcastHeartbeatLocked()
				}
			} else if !now.Before(rf.electionDeadline) {
				rf.startElectionLocked()
			}
			rf.mu.Unlock()
		}
	}
}

// Start preserva a assinatura do esqueleto. O indice -1 indica que nenhum
// comando foi anexado: replicacao e ApplyMsg pertencem a parte 2B.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	term, isLeader := rf.GetState()
	return -1, term, isLeader
}

func (rf *Raft) Kill() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if !rf.stopped {
		rf.stopped = true
		rf.state = follower
		close(rf.stopCh)
	}
}

func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{
		peers: peers, me: me, persister: persister,
		votedFor: -1, state: follower,
		log: []LogEntry{{Term: 0}},
		rng: rand.New(rand.NewSource(time.Now().UnixNano() + int64(me)*1000003)),
		stopCh: make(chan struct{}),
	}
	// Nenhum arquivo e acessado. Persister e applyCh ficam na API para
	// compatibilidade; persistencia e aplicacao nao sao usadas em 2A.
	rf.resetElectionLocked()
	go rf.ticker()
	return rf
}
