"""Cria um ZIP somente com fontes, scripts e documentacao do trabalho."""
from pathlib import Path
from zipfile import ZIP_DEFLATED, ZipFile

root = Path(__file__).resolve().parent.parent
output = root / 'output' / 'trabalho-raft-2a.zip'
output.parent.mkdir(exist_ok=True)
files = [root / 'README.md', root / '.gitignore']
for folder in ('src', 'docs', 'scripts'):
    files.extend(p for p in (root / folder).rglob('*')
                 if p.is_file() and p.suffix in ('.go', '.md', '.json', '.ps1', '.sh', '.py'))
with ZipFile(output, 'w', ZIP_DEFLATED) as archive:
    for path in sorted(files):
        archive.write(path, path.relative_to(root).as_posix())
with ZipFile(output) as archive:
    if archive.testzip() is not None:
        raise RuntimeError('Falha ao verificar o ZIP')
    print(f'ZIP verificado: {len(archive.namelist())} arquivos; {output.stat().st_size} bytes')
print(output)
