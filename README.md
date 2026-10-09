# zero-kb02 Codex controller

[English](README_EN.md)

zero-kb02を、macOSのHerdrで動くCodex CLI用のコントローラーにするプロジェクトです。最大6つの会話の状態をOLEDとLEDに表示し、キーで会話を切り替え、Encoderで推論の強さを変更できます。Joystickはマウスポインターの移動に使えます。

[デモ動画（X）](https://x.com/hoki621/status/2093605017819423047)

## システム構成

![構成図](assets/architecture-ja.drawio.svg)

zero-kb02上のFirmwareが入力と表示を担当します。Mac上のHost bridgeがデバイスとHerdrを接続し、Codex App Serverを通じて会話の推論設定を変更します。

## 用意するもの

- zero-kb02とUSBデータケーブル
- macOSと[Herdr](https://herdr.dev/)
- [Homebrew](https://brew.sh/)、[mise](https://mise.jdx.dev/)、Codex CLIを利用できるアカウント

Node.js・Go・TinyGoは`mise.toml`でバージョンを指定しています。Codex CLIはHomebrew版を使用します。

## 初期設定

1. Herdr、Homebrew、miseをインストールし、次のコマンドでソースと必要なツールを取得します。

   ```sh
   brew install --cask codex
   git clone --recurse-submodules https://github.com/hoki621/codex-zero-kb02.git
   cd codex-zero-kb02
   mise install
   mise exec -- sh -c 'cd host && npm ci && npm run build'
   ```

2. Codex CLIの認証と、Herdrの状態一覧pluginの登録を行います。

   ```sh
   codex login
   herdr plugin link --enabled "$PWD/host"
   ```

3. [Firmwareの手順](firmware/README_JA.md)に沿ってビルド・書き込みを行い、zero-kb02をUSBで接続します。

## 起動

Herdrを起動してから、以下をリポジトリのルートで実行します。それぞれ別のTerminalまたはペインを使ってください。

1. 通常のTerminalでCodex App Serverを起動し、そのままにします。会話の実行と推論設定を管理するサービスです。

   ```sh
   mise exec -- node host/dist/src/codex-micro.js server
   ```

2. Herdrのペインで、デバイスと連携するCodex CLIを起動します。会話ごとに1つのペインを使います（最大6つ）。

   ```sh
   mise exec -- node host/dist/src/codex-micro.js
   ```

   この起動用プログラムがCodexの会話とHerdrのペインを結びつけます。推論変更・承認操作を使う会話は、通常の`codex`ではなくこのコマンドで起動してください。会話の再開は末尾に`resume`を付けます。別のプロジェクトで使う場合は、そのディレクトリへ移動し、`host/dist/src/codex-micro.js`を絶対パスで指定します。

3. 別の通常TerminalでHost bridgeを起動します。USBポート名は自分のデバイスのものに置き換え、ほかのserial monitorは閉じてください。ポートの候補は`ls /dev/cu.usbmodem*`で確認できます。

   ```sh
   HERDR_SOCKET_PATH="$HOME/.config/herdr/herdr.sock" \
   ZERO_KB02_PORT=/dev/cu.usbmodemzero_kb02_v21 \
   mise exec -- node host/dist/src/main.js
   ```

   `HERDR_SOCKET_PATH`はHerdrとの接続先、`ZERO_KB02_PORT`はzero-kb02のUSB接続先です。

終了するときは、各Codex CLIとHost bridgeを終了してから、App ServerのTerminalでCtrl-Cを押します。

## 操作

キー番号は上段左からK1〜K4、中段K5〜K8、下段K9〜K12です。

| 入力 | 動作 |
| --- | --- |
| K1 | 選択中のCodexへEscape |
| K2 / K3 / K5 / K6 / K7 / K8 | 会話1〜6のペインへ切替 |
| K4 | 状態一覧を開く・閉じる |
| K9 / K10 | コマンド実行の承認・拒否（対応版のみ） |
| K11 | 未割当 |
| K12 | 入力待ち・完了のCodexで新しい会話 |
| Encoder | 時計回りで推論の強さを上げ、反時計回りで下げる |
| Joystick | マウスポインターの移動。押し込み操作は未割当 |

OLEDは会話1/2、3/4、5/6の2列3段です。W=作業中、I=入力待ち、B=承認・回答待ち、D=完了、U=不明、E=空枠を表します。

K9/K10はCodex CLI **0.155.1・0.160.0のみ対応**し、操作対象と承認画面を確認できる場合に有効です。**0.162.0では無効**です。K4はHerdr共通のpopupを操作するため、別pluginのpopupを閉じる場合があります。

## 更新・開発

更新時はCLI・Host bridge・App Serverを終了してから、リポジトリのルートで実行します。

```sh
git pull --ff-only
git submodule update --init --recursive
mise install
mise exec -- sh -c 'cd host && npm ci && npm run build'
```

Codex CLIをHomebrewで更新した場合も、App Serverと各CLIを再起動してください。ソース変更の確認は実機なしで行えます。

```sh
mise exec -- sh -c 'cd host && npm run typecheck && npm test && npm run dry-run -- WIBDUE'
mise exec -- sh -c 'cd firmware && go test ./... && go vet ./...'
```

- [Hostの詳細・トラブルシューティング](host/README_JA.md) / [Firmware](firmware/README_JA.md)
- [動作確認の範囲](docs/verification.md) / [通信仕様](PROTOCOL.md)
- [使用ライブラリとライセンス](UPSTREAMS.md)

上流ライブラリのライセンスは各子リポジトリに保持しています。このプロジェクト独自のコードにはまだライセンスを指定していません。
