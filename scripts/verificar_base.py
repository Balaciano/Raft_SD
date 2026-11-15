"""Verifica os arquivos fornecidos usando o hash de blob do Git (sem rede)."""
import hashlib
import json
from pathlib import Path

root = Path(__file__).resolve().parent.parent
manifest = json.loads((root / 'docs/base-original.json').read_text(encoding='utf-8'))
failed = []
for item in manifest['arquivos']:
    data = (root / item['path']).read_bytes()
    digest = hashlib.sha1(b'blob ' + str(len(data)).encode() + b'\0' + data).hexdigest()
    if digest != item['git_blob_sha1']:
        failed.append(item['path'])
        print('DIFERENTE:', item['path'])
    else:
        print('OK:', item['path'])
raise SystemExit(1 if failed else 0)
