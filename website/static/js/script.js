var logsTable = document.getElementById("logsTable");

window.onload = function () {
  setInterval(function () {
    getLogs();
  }, 1000);
};

async function getLogs() {
  fetch("http://localhost:8412/logs")
    .then(function (response) {
      return response.json();
    })
    .then(function (data) {
      clearTable();
      generateTable(data.logs);
    });
}

function clearTable() {
  logsTable.innerHTML = "";
}

function generateTable(logs) {
  for (let i = 0; i < logs.length; i++) {
    const pre = document.createElement("pre");
    const code = document.createElement("code");
    code.classList = "language-json";

    const cellText = document.createTextNode(`${logs[i]}`);

    code.appendChild(cellText);
    pre.appendChild(code);

    logsTable.appendChild(pre);

    hljs.highlightAll();
  }
}
