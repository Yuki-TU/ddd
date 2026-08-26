#!/usr/bin/env python3
"""diagrams/*.puml を SVG に変換する。初回は README の ```plantuml も画像参照に置き換える。"""

from __future__ import annotations

import re
import shutil
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
README = ROOT / "README.md"
DIAGRAMS = ROOT / "diagrams"

NAMES = [
    ("ch06-usecase", "ユーザ機能のクラス図"),
    ("ch06-register-user", "ユーザ登録"),
    ("ch06-get-user", "ユーザ情報の取得"),
    ("ch06-update-user", "ユーザ情報の更新"),
    ("ch06-withdraw", "退会処理"),
    ("ch06-cohesion", "ユースケースごとのアプリケーションサービス"),
    ("ch06-package", "アプリケーションサービスのパッケージ構成"),
    ("ch07-dependency", "依存関係"),
    ("ch07-interface", "インターフェースの実現"),
    ("ch07-multiple-impl", "複数のリポジトリ実装"),
    ("ch07-dip-violation", "依存関係逆転の原則への違反"),
    ("ch07-dip", "依存関係逆転の原則"),
    ("ch07-interface-ownership", "インターフェースの主導権"),
    ("ch08-cli", "CLI への組み込み"),
    ("ch08-cli-sequence", "CLI のシーケンス"),
    ("ch08-startup", "サーバー起動時の依存関係登録"),
    ("ch08-mvc", "コントローラーとアプリケーションサービス"),
    ("ch08-mvc-sequence", "ユーザ登録リクエストのシーケンス"),
    ("ch11-usecase", "サークルのユースケース"),
    ("ch11-circle-create", "サークル作成の構成"),
    ("ch12-aggregate-user", "User 集約"),
    ("ch12-aggregate-circle", "Circle 集約"),
    ("ch12-aggregate-userid", "Circle 集約に UserId だけ含める"),
    ("ch13-specification", "仕様とリポジトリ"),
]


def wrap_puml(source: str) -> str:
    body = source.strip()
    if not body.startswith("@startuml"):
        body = "@startuml\n" + body
    if not body.rstrip().endswith("@enduml"):
        body = body.rstrip() + "\n@enduml\n"
    return body


def render_all() -> int:
    plantuml = shutil.which("plantuml")
    if not plantuml:
        print("plantuml が見つかりません。`brew install plantuml graphviz` を実行してください。", file=sys.stderr)
        return 1
    puml_files = sorted(DIAGRAMS.glob("*.puml"))
    if not puml_files:
        print("diagrams/*.puml がありません。", file=sys.stderr)
        return 1
    result = subprocess.run(
        [plantuml, "-tsvg", "-charset", "UTF-8", *[str(p) for p in puml_files]],
        check=False,
        capture_output=True,
        text=True,
    )
    if result.stdout:
        print(result.stdout)
    if result.stderr:
        print(result.stderr, file=sys.stderr)
    missing = [p.with_suffix(".svg").name for p in puml_files if not p.with_suffix(".svg").exists()]
    if result.returncode != 0 or missing:
        if missing:
            print("未生成: " + ", ".join(missing), file=sys.stderr)
        return 1
    for p in puml_files:
        svg = p.with_suffix(".svg")
        print(f"ok  {svg.name} ({svg.stat().st_size} bytes)")
    return 0


def replace_readme_fences() -> int:
    text = README.read_text(encoding="utf-8")
    pattern = re.compile(r"```plantuml\n(.*?)```", re.DOTALL)
    blocks = list(pattern.finditer(text))
    if not blocks:
        return 0
    if len(blocks) != len(NAMES):
        print(f"block count {len(blocks)} != name count {len(NAMES)}", file=sys.stderr)
        return 1

    DIAGRAMS.mkdir(exist_ok=True)
    out_parts: list[str] = []
    last = 0
    for i, match in enumerate(blocks):
        slug, title = NAMES[i]
        puml = wrap_puml(match.group(1))
        (DIAGRAMS / f"{slug}.puml").write_text(puml, encoding="utf-8")
        out_parts.append(text[last : match.start()])
        out_parts.append(f"![{title}](./diagrams/{slug}.svg)")
        last = match.end()
    out_parts.append(text[last:])
    README.write_text("".join(out_parts), encoding="utf-8")
    return 0


def main() -> int:
    if replace_readme_fences() != 0:
        return 1
    return render_all()


if __name__ == "__main__":
    raise SystemExit(main())
