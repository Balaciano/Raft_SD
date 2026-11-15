package raft

import (
	"labrpc"
	"math/rand"
	"testing"
	"time"
)

// Peer sem ticker para verificar transicoes sem depender do escalonador.
func isolatedPeer() *Raft {
	return &Raft{
		currentTerm: 3, votedFor: -1, state: follower,
		log: []LogEntry{{Term: 0}},
		rng: rand.New(rand.NewSource(1)), stopCh: make(chan struct{}),
		electionDeadline: time.Now().Add(time.Second),
	}
}

func TestOneVotePerTerm2A(t *testing.T) {
	rf := isolatedPeer()
	var reply RequestVoteReply
	rf.RequestVote(&RequestVoteArgs{Term: 3, CandidateId: 1}, &reply)
	if !reply.VoteGranted {
		t.Fatal("primeiro voto foi recusado")
	}
	rf.RequestVote(&RequestVoteArgs{Term: 3, CandidateId: 2}, &reply)
	if reply.VoteGranted {
		t.Fatal("dois candidatos receberam voto no mesmo termo")
	}
	rf.RequestVote(&RequestVoteArgs{Term: 3, CandidateId: 1}, &reply)
	if !reply.VoteGranted {
		t.Fatal("repeticao idempotente do pedido foi recusada")
	}
	rf.RequestVote(&RequestVoteArgs{Term: 4, CandidateId: 2}, &reply)
	if !reply.VoteGranted || reply.Term != 4 {
		t.Fatal("novo termo nao liberou o voto")
	}
}

func TestStaleRPCDoesNotResetTimer2A(t *testing.T) {
	rf := isolatedPeer()
	before := rf.electionDeadline
	var vote RequestVoteReply
	rf.RequestVote(&RequestVoteArgs{Term: 2, CandidateId: 1}, &vote)
	var appendReply AppendEntriesReply
	rf.AppendEntries(&AppendEntriesArgs{Term: 2}, &appendReply)
	if vote.VoteGranted || appendReply.Success || !rf.electionDeadline.Equal(before) {
		t.Fatal("RPC obsoleta foi aceita ou adiou a eleicao")
	}
}

func TestDeniedVoteDoesNotResetTimer2A(t *testing.T) {
	rf := isolatedPeer()
	rf.log = append(rf.log, LogEntry{Term: 3})
	before := rf.electionDeadline
	var reply RequestVoteReply
	rf.RequestVote(&RequestVoteArgs{Term: 4, CandidateId: 1}, &reply)
	if reply.VoteGranted || reply.Term != 4 || !rf.electionDeadline.Equal(before) {
		t.Fatal("pedido com log atrasado deve atualizar termo sem reiniciar timer")
	}
}

func TestSameTermHeartbeatPreservesVote2A(t *testing.T) {
	rf := isolatedPeer()
	rf.state, rf.votedFor = candidate, 0
	var reply AppendEntriesReply
	rf.AppendEntries(&AppendEntriesArgs{Term: 3, LeaderId: 1}, &reply)
	if !reply.Success || rf.state != follower || rf.votedFor != 0 {
		t.Fatal("heartbeat deve encerrar candidatura sem apagar voto do termo")
	}
	var vote RequestVoteReply
	rf.RequestVote(&RequestVoteArgs{Term: 3, CandidateId: 2}, &vote)
	if vote.VoteGranted {
		t.Fatal("heartbeat permitiu votar novamente no mesmo termo")
	}
}

func TestHigherTermStepsDown2A(t *testing.T) {
	rf := isolatedPeer()
	rf.state, rf.votedFor = leader, 0
	var reply AppendEntriesReply
	rf.AppendEntries(&AppendEntriesArgs{Term: 6, LeaderId: 1}, &reply)
	if !reply.Success || rf.state != follower || rf.currentTerm != 6 || rf.votedFor != -1 {
		t.Fatal("lider nao reconheceu termo superior")
	}
}

func TestHeartbeatChecksPrefix2A(t *testing.T) {
	rf := isolatedPeer()
	var reply AppendEntriesReply
	rf.AppendEntries(&AppendEntriesArgs{Term: 3, PrevLogIndex: 1}, &reply)
	if reply.Success {
		t.Fatal("heartbeat aceitou prefixo inexistente")
	}
	rf.AppendEntries(&AppendEntriesArgs{Term: 3, PrevLogTerm: 8}, &reply)
	if reply.Success {
		t.Fatal("heartbeat aceitou termo de prefixo diferente")
	}
}

func TestKillIsIdempotent2A(t *testing.T) {
	rf := isolatedPeer()
	rf.Kill()
	rf.Kill()
	var reply RequestVoteReply
	rf.RequestVote(&RequestVoteArgs{Term: 9, CandidateId: 1}, &reply)
	term, isLeader := rf.GetState()
	if reply.VoteGranted || isLeader || term != 3 {
		t.Fatal("peer encerrado continuou participando")
	}
}

func TestSingleServerElection2A(t *testing.T) {
	rf := Make(make([]*labrpc.ClientEnd, 1), 0, MakePersister(), make(chan ApplyMsg))
	defer rf.Kill()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if term, isLeader := rf.GetState(); isLeader && term == 1 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("servidor unico nao se elegeu com seu proprio voto")
}

func TestUnreliableElection2A(t *testing.T) {
	cfg := make_config(t, 5, true)
	defer cfg.cleanup()
	leader := cfg.checkOneLeader()
	cfg.disconnect(leader)
	cfg.checkOneLeader()
	cfg.connect(leader)
	cfg.checkOneLeader()
}

func TestReelectionTiming2A(t *testing.T) {
	cfg := make_config(t, 3, false)
	defer cfg.cleanup()
	oldLeader := cfg.checkOneLeader()
	cfg.disconnect(oldLeader)
	started := time.Now()
	for time.Since(started) < 5*time.Second {
		for i := 0; i < cfg.n; i++ {
			if cfg.connected[i] {
				if _, isLeader := cfg.rafts[i].GetState(); isLeader {
					t.Logf("REELEICAO_SEGUNDOS=%.3f (rede confiavel; resolucao de consulta 10 ms)", time.Since(started).Seconds())
					return
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("nao houve reeleicao dentro de cinco segundos com maioria conectada")
}

func TestHeartbeatRate2A(t *testing.T) {
	cfg := make_config(t, 3, false)
	defer cfg.cleanup()
	leader := cfg.checkOneLeader()
	time.Sleep(500 * time.Millisecond)
	term, isLeader := cfg.rafts[leader].GetState()
	if !isLeader {
		t.Fatal("lider mudou antes da medicao")
	}
	target := (leader + 1) % cfg.n
	before := cfg.rpcCount(target)
	started := time.Now()
	time.Sleep(2 * time.Second)
	count := cfg.rpcCount(target) - before
	elapsed := time.Since(started).Seconds()
	afterTerm, stillLeader := cfg.rafts[leader].GetState()
	if !stillLeader || afterTerm != term || cfg.checkTerms() != term {
		t.Fatal("medicao de heartbeats invalidada por mudanca de termo/lider")
	}
	rate := float64(count) / elapsed
	t.Logf("HEARTBEATS_RPC=%d JANELA_SEGUNDOS=%.3f TAXA_RPC_POR_SEGUNDO=%.3f", count, elapsed, rate)
	if count == 0 || rate > 10 {
		t.Fatalf("taxa observada fora do limite: %.3f RPCs/s", rate)
	}
}
