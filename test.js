console.log("✅ JS runtime started!");

setTimeout(() => {
    console.log("⏰ This is printed after 2 seconds!");
}, 2000);

const content = fs.readFileSync("hello.txt");
console.log("📄 File Content:", content);

http.createServer((req, res) => {
    console.log("🌐 Received request for:", req.url);
    res.end("Hello from Go-powered JS server!");
});
