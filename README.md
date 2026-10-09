# zero-kb02 Codex controller

zero-kb02（RP2040）の12キー・Encoder・OLED/LEDを使い、macOSのHerdr上で動く最大6つのCodex CLIを表示・操作します。
Codex CLIはHomebrewで管理します。TinyGo Firmwareは固定した公開ライブラリを組み合わせて実装します。

**PC側は実装・独立レビュー済みです。Firmware major 2は[子repoのPR #1](https://github.com/hoki621/codex-zero-kb02-firmware/pull/1)でビルド済み、実機受入は未完了です。**
HostはUSB **major 2**を使用します。親が保持する既存Firmware `4d8104c`はmajor 1なので接続できません。
[段階別のFirmware手順](docs/firmware-handoff.md)は実機受入の順序として使います。親gitlinkは実機受入後に更新します。

## 構成と操作

- `host/`: Node.js 22。Herdrの状態取得、6枠管理、CDC通信、固定キー操作、Codex App Server接続、開発CLI。
- `firmware/`: TinyGo実装。matrix/debounce、Encoder、OLED、LEDは公開ライブラリ、Joystickは標準HID mouse。製品固有のCDC契約と表示対応を接続します。
- [PROTOCOL.md](PROTOCOL.md): major 2の通信契約。[UPSTREAMS.md](UPSTREAMS.md): 固定依存とライセンス。
- [計画 #34](https://github.com/hoki621/codex-zero-kb02/issues/34): 作業管理。[検証記録](docs/verification.md): mock/API/buildと実機の区別。

```text
┌────────┬────────┬────────┬────────┐
│ K1     │ K2     │ K3     │ K4     │
│ Escape │ Agent1 │ Agent2 │ Status │
├────────┼────────┼────────┼────────┤
│ K5     │ K6     │ K7     │ K8     │
│ Agent3 │ Agent4 │ Agent5 │ Agent6 │
├────────┼────────┼────────┼────────┤
│ K9     │ K10    │ K11    │ K12    │
│ Approve│ Reject │ 未割当 │ /new   │
└────────┴────────┴────────┴────────┘
```

Encoderの時計回りはreasoning effortを1クリック1段階上げ、反時計回りは下げます。
Joystickはポインター移動のみ。両pushは未割当です。
OLEDはAgent1/2、3/4、5/6を2列3段に並べます。W=作業中、I=入力待ち、B=承認・回答待ち、D=完了、U=不明、E=空枠です。

K9/K10は、対応版の単独command承認・会話ID・terminal・画面を確認できるときだけ固定y/nを各1回送ります。
`blocked`だけでは送信しません。質問・未知の画面・複数の承認は無効です。永続承認やEnterは送りません。
K12はfocus中の割当済みCodexがidle/doneの場合だけ固定/newを送ります。

## キーボードなしで準備・確認

```sh
git clone --recurse-submodules https://github.com/hoki621/codex-zero-kb02.git
cd codex-zero-kb02
# mise.tomlに指定された版が未導入の場合だけmise install
cd host
npm ci
npm run typecheck
npm test
npm run device:check -- input
npm run device:check -- display
npm run device:check -- faults
npm run dry-run -- WIBDUE
cd ..
sh docs/library-probe/check.sh
```

device:checkはmockが既定です。`raw`は初期入力ログ用、`doctor`は実行ファイル・版・socket・通信majorの診断です。
依存probeは一時ディレクトリでビルドするだけで、書き込みません。64のHostテストと独立レビューの詳細は[検証記録](docs/verification.md)へ。

brew版App Serverを試す場合:

```sh
cd host
npm run smoke:codex
```

これは一時CODEX_HOMEのephemeral会話で、最初の発言前にreasoningを1段階変更します。
モデルへの発言・実Herdrへの入力・USB接続はなく、終了時に専用プロセスと一時領域を削除します。

## Firmware完成後の起動

最初に[段階別手順](docs/firmware-handoff.md)のH0〜H4で、major 2対応Firmwareを実機確認してください。
対象portは完全な名前を指定します。自動探索はしません。開発用CLI・他のserial monitorを終了してからbridgeを起動します。

初回だけHost実行ファイルと固定Status pluginを登録します:

```sh
cd host
npm run build
npm link
herdr plugin link --enabled "$(pwd)"
```

1. 通常Terminalでbrew版の専用serverを起動したままにします。

   ```sh
   brew install --cask codex  # 未導入の場合だけ
   codex-micro doctor
   codex-micro server
   ```

2. Herdrの各Codex paneで`codex-micro`を実行します。再開は`codex-micro resume`です。
3. 別の通常TerminalからHost bridgeを起動します。Herdrの管理pane内には置きません。

   ```sh
   cd host
   HERDR_SOCKET_PATH="$HOME/.config/herdr/herdr.sock" \
   ZERO_KB02_PORT=/dev/cu.usbmodemYOUR_DEVICE \
   npm start
   ```

CLIとserverは同じbrew管理バイナリを使用します。`codex-micro`はstart/resume/forkの応答から正確な会話UUIDv7を登録します。
Host bridgeを止めても専用serverや他の会話は止まりません。終了は各TerminalのCtrl-Cです。
serverは使用中のremote CLIをすべて終了してから止めます。

## 復旧と更新

- USB抜き差し: bridgeは同じportへ再接続し、全snapshotを再送します。再接続時に押していたキーは離してから操作します。
- Herdr再起動: bridgeが再接続します。5秒ごとの再照合で状態を修復します。
- Codexのbrew更新: remote CLIを終了し、専用serverをCtrl-Cで停止して再起動し、各paneでresumeします。対応外の承認画面は無効になります。
- serial使用中: `lsof /exact/port`で所有者を確認し、自分のbridge/monitorを終了します。不明なプロセスを強制終了しません。
- server/launcherの異常終了: [Host復旧手順](host/README.md#pc-setup-and-startup)でPID/socketを確認します。live ownerのファイルを削除しません。
- Firmwareの復旧: [実機診断手順](firmware/docs/hardware-diagnostics.md)。Codexに実行させる場合、flash・port openは対象操作を明示して依頼してください。

親commitに固定したソースを取得します。`git submodule update --remote`は使いません。

```sh
git pull --ff-only
git submodule update --init --recursive
# mise.tomlに指定された版が未導入の場合だけmise install
cd host
npm ci
```

過去のSessionStart hookを使用していた場合は[Hostの旧hook移行手順](host/README.md#legacy-hook-migration)も参照してください。新規導入にはhook不要です。

## 対応範囲と未検証項目

Herdr 0.9.0のJSON schemaを参照し、必須フィールドを検証しています。内部protocol番号だけでは拒否しません。
Codex 0.155.1と0.160.0の隔離API試験が成功しています。現在の親gitlinkが参照するHostは承認対応版を0.155.1に限定します。[Host PR #3](https://github.com/hoki621/codex-zero-kb02-host/pull/3)と[親PR #38](https://github.com/hoki621/codex-zero-kb02/pull/38)を適用すると0.160.0も対応し、process-onlyのkeymap指定でy/nを固定します。
設定ファイルは変更しません。K4のpopupはHerdr session共通で、別pluginのpopupを閉じる場合があります。

実Herdrの対話画面、USB列挙、物理入力、OLED/LED、抜き差し、flashは**NOT RUN**です。
[実機受入 #32](https://github.com/hoki621/codex-zero-kb02/issues/32)で確認します。PC試験の成功を実機成功として扱いません。
launchd、自動起動、設定GUI、Vial制御、任意入力/任意shell、永続承認、PTT、model切替、Codex Desktop/Zed対応は範囲外です。

開発は親Issueで管理し、childを先にcommit/pushしてから親gitlinkを更新します。Firmwareの実機受入はコードのビルド結果と分けて記録します。
