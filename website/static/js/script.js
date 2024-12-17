var logsTable = document.getElementById("logsTable");
var logsAmount = document.getElementById("logsAmount");
var logsAmountInfo = document.getElementById("logsAmountInfo");
var logsAmountWarning = document.getElementById("logsAmountWarning");
var logsAmountError = document.getElementById("logsAmountError");

var pause = false;
var filtered = undefined;
var logEntries = {
  INFO: {},
  WARNING: {},
  ERROR: {},
  DEBUG: {},
};

window.onload = function () {
  setInterval(function () {
    if (!pause) {
      getLogs();
    }
  }, 1000);
};

function toggleLogs() {
  pause = !pause;
}

function updateTopBar(amountOfLogs, counts) {
  logsAmount.innerHTML = "logs " + amountOfLogs;
  logsAmountInfo.innerHTML = "info " + counts.INFO;
  logsAmountWarning.innerHTML = "warn " + counts.WARNING;
  logsAmountError.innerHTML = "err " + counts.ERROR;
}

function getLogs() {
  fetch("http://localhost:8412/logs")
    .then(function (response) {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then(function (data) {
      clearTable();
      generateTable(data.logs);
    })
    .catch(function (error) {
      showNotification("Failed to fetch logs");
      console.error("Failed to fetch logs:", error);
    });
}

function clearTable() {
  logsTable.innerHTML = "";
}

function generateTable(logs) {
  const levelMap = {
    info: "INFO",
    warn: "WARNING",
    warning: "WARNING",
    error: "ERROR",
    err: "ERROR",
    debug: "DEBUG",
  };

  let counts = {
    INFO: 0,
    WARNING: 0,
    ERROR: 0,
    DEBUG: 0,
  };

  logs.forEach((logEntry, index) => {
    let parsedLog;
    let logLevelClass = "DEFAULT";

    try {
      parsedLog = JSON.parse(logEntry).level.toLowerCase();
      if (levelMap[parsedLog]) {
        logLevelClass = levelMap[parsedLog];
        counts[logLevelClass]++;
      }
    } catch (e) {
      console.error("Failed to parse log entry:", logEntry, e);
    }

    const pre = document.createElement("pre");
    const code = document.createElement("code");
    const cellText = document.createTextNode(logEntry);

    pre.id = "log-id-" + index;
    pre.classList.add(logLevelClass);
    code.classList.add("language-json");

    code.appendChild(cellText);
    pre.appendChild(code);

    showLogEntry(pre);
    logEntries[logLevelClass][pre.id] = pre;
  });

  updateTopBar(logs.length, counts);
  hljs.highlightAll();
}

function showLogEntry(logEntry) {
  if (filtered) {
    if (logEntry.classList.contains(filtered.toUpperCase())) {
      logsTable.appendChild(logEntry);
    }
  } else {
    logsTable.appendChild(logEntry);
  }
}

function showNotification(message) {
  const notification = document.createElement("div");
  notification.classList.add("notification");

  notification.textContent = message;
  document.body.appendChild(notification);

  setTimeout(() => {
    notification.style.opacity = "0";
  }, 5000);

  setTimeout(() => {
    notification.remove();
  }, 6000);
}

function filterLogs(level) {
  filtered = level;
  for (let logLevel in logEntries) {
    if (logLevel === level.toUpperCase()) continue;

    for (let elementId in logEntries[logLevel]) {
      var logElement = logEntries[logLevel][elementId];
      logElement.style.visibility = "hidden";
    }
  }
}
