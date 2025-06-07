
# ⚡ VoltJS

A lightweight, Node.js-inspired JavaScript runtime powered by [Goja](https://github.com/dop251/goja) and written in Go. VoltJS provides native support for `console`, `setTimeout`, `fs`, HTTP servers, and a simple `require` system — ideal for scripting and server-side JavaScript powered by Go's performance.

---

## 🚀 Features

- 🧠 **JavaScript runtime** with CommonJS-like `require` system
- 🕒 `setTimeout`, `setInterval`, `clearTimeout`, `clearInterval`
- 📂 `fs.readFileSync`, `fs.writeFileSync`
- 🌐 Simple `http.createServer` API
- 📦 `process.argv`, `process.cwd()`
- 🔀 Module isolation via new `goja.Runtime` per module
- 🔧 CLI-friendly and hackable

---

## 🧪 Example Usage

### `script.js`
```js
console.log("⚡ Running from VoltJS!");

setTimeout(() => {
    console.log("✅ Done after 1 second!");
}, 1000);

const data = fs.readFileSync("hello.txt");
console.log("📄 File content:", data);

http.createServer((req, res) => {
    console.log("📥 Incoming request:", req.url);
    res.end("Hello from VoltJS server!");
});
````

### Run:

```bash
voltjs script.js
```

---

## 🛠 Installation

You can install the latest prebuilt binary using the following one-liner:

```bash
curl -s https://raw.githubusercontent.com/Pjdeveloper896/VoltJs/main/install.sh | bash
```

Or manually:

### 📜 install.sh

```bash
#!/usr/bin/env bash

set -e

REPO="Pjdeveloper896/VoltJs"
BIN_NAME="voltjs"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "❌ Unsupported architecture: $ARCH"
  exit 1
fi

if [[ -n "$PREFIX" ]]; then
  INSTALL_DIR="$PREFIX/bin"
elif [[ $EUID -ne 0 ]]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
else
  INSTALL_DIR="/usr/local/bin"
fi

echo "📦 Installing to $INSTALL_DIR ..."

LATEST_RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")

TAG_NAME=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"tag_name": "\K(.*)(?=")')
if [[ -z "$TAG_NAME" ]]; then
  echo "❌ Could not fetch the latest release tag."
  exit 1
fi

echo "🔖 Latest release: $TAG_NAME"

FILENAME="${BIN_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"browser_download_url": "\K(.*)(?=")' | grep "$FILENAME")

if [[ -z "$DOWNLOAD_URL" ]]; then
  echo "❌ No binary found for ${FILENAME} in latest release."
  exit 1
fi

echo "⬇️  Downloading $FILENAME from:"
echo "    $DOWNLOAD_URL"

TMP_FILE="/tmp/$FILENAME"
curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"

chmod +x "$TMP_FILE"
mv "$TMP_FILE" "$INSTALL_DIR/$BIN_NAME"

echo "✅ Installed $BIN_NAME to $INSTALL_DIR"

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "⚠️  $INSTALL_DIR is not in your PATH."
  echo "➡️  Add this to your shell config:"
  echo "    export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "🚀 Run it with:"
echo "    $BIN_NAME --help"
```

---

## 📁 File Structure

```
VoltJs/
├── main.go            # VoltJS runtime implementation in Go
├── install.sh         # One-line installer script
├── example/
│   └── script.js      # Sample JS script
├── modules/           # Custom JS modules loaded with require()
└── README.md
```

---

## 🔧 Build from Source

```bash
git clone https://github.com/Pjdeveloper896/VoltJs.git
cd VoltJs
go build -o voltjs main.go
./voltjs example/script.js
```

---

## 🧩 Future Plans

* [ ] `fetch()` API support via Go HTTP client
* [ ] ES module system (ESM)
* [ ] Support for Promises
* [ ] Native `__dirname` and `__filename` support
* [ ] More built-in modules: `os`, `path`, etc.

---

## 🧑‍💻 Author

**Prasoon Jadon**
📦 GitHub: [@Pjdeveloper896](https://github.com/Pjdeveloper896)

---

