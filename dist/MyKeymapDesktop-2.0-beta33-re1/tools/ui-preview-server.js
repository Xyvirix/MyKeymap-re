const fs = require("fs");
const http = require("http");
const path = require("path");
const rootDir = path.resolve(__dirname, "..");
const siteDir = path.join(rootDir, "config-ui", "dist");
const configFile = path.join(rootDir, "data", "config.json");
const shortcutsDir = path.join(rootDir, "shortcuts");
const port = Number(process.env.UI_PREVIEW_PORT || 4173);

const contentTypes = {
  ".css": "text/css; charset=utf-8",
  ".eot": "application/vnd.ms-fontobject",
  ".gif": "image/gif",
  ".html": "text/html; charset=utf-8",
  ".ico": "image/x-icon",
  ".jpg": "image/jpeg",
  ".js": "text/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".md": "text/markdown; charset=utf-8",
  ".png": "image/png",
  ".svg": "image/svg+xml",
  ".ttf": "font/ttf",
  ".txt": "text/plain; charset=utf-8",
  ".woff": "font/woff",
  ".woff2": "font/woff2",
};

function sendJson(res, status, payload) {
  res.writeHead(status, { "Content-Type": "application/json; charset=utf-8" });
  res.end(JSON.stringify(payload));
}

function sendFile(res, filePath) {
  fs.readFile(filePath, (err, data) => {
    if (err) {
      res.writeHead(404, { "Content-Type": "text/plain; charset=utf-8" });
      res.end("Not found");
      return;
    }

    const ext = path.extname(filePath).toLowerCase();
    res.writeHead(200, {
      "Content-Type": contentTypes[ext] || "application/octet-stream",
      "Cache-Control": "no-store",
    });
    res.end(data);
  });
}

function listShortcuts() {
  if (!fs.existsSync(shortcutsDir)) {
    return [];
  }

  return fs.readdirSync(shortcutsDir)
    .filter(name => name.toLowerCase().endsWith(".lnk"))
    .map(name => ({ path: `shortcuts\\${name}` }));
}

function resolveSiteFile(requestPath) {
  const cleanPath = requestPath === "/" ? "/index.html" : requestPath;
  const normalized = path.normalize(cleanPath).replace(/^(\.\.[/\\])+/, "");
  let filePath = path.join(siteDir, normalized);

  if (fs.existsSync(filePath) && fs.statSync(filePath).isFile()) {
    return filePath;
  }

  filePath = path.join(siteDir, "index.html");
  return filePath;
}

const server = http.createServer((req, res) => {
  const parsed = new URL(req.url || "/", `http://${req.headers.host || "127.0.0.1"}`);
  const requestPath = parsed.pathname || "/";

  if (requestPath === "/config" && req.method === "GET") {
    return sendFile(res, configFile);
  }

  if (requestPath === "/config" && req.method === "PUT") {
    req.resume();
    req.on("end", () => sendJson(res, 200, { message: "ok" }));
    return;
  }

  if (requestPath === "/shortcuts" && req.method === "GET") {
    return sendJson(res, 200, listShortcuts());
  }

  if (requestPath.startsWith("/server/command/") && req.method === "POST") {
    req.resume();
    req.on("end", () => sendJson(res, 200, {}));
    return;
  }

  return sendFile(res, resolveSiteFile(requestPath));
});

server.listen(port, "127.0.0.1", () => {
  console.log(`UI preview server listening on http://127.0.0.1:${port}`);
});
