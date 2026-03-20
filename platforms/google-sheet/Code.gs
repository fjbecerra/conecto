function onOpen() {
  SpreadsheetApp.getUi()
    .createMenu('HTTP Tools')
    .addItem('Open HTTP Request', 'showSidebar')
    .addToUi();
}

function showSidebar() {
  const html = HtmlService.createHtmlOutputFromFile('sidebar')
    .setWidth(350);
  SpreadsheetApp.getUi().showModelessDialog(html, 'HTTP Request');
}

// ADD THIS - Handle the HTTP request from backend
function makeHttpRequest(method, url) {
  try {
    const encodedUrl = encodeURI(url);
    const options = {
      method: method,
      muteHttpExceptions: true
    };
    
    const response = UrlFetchApp.fetch(encodedUrl, options);
    const responseCode = response.getResponseCode();
    const responseText = response.getContentText();
    
    return {
      success: responseCode === 200,
      status: responseCode,
      data: JSON.parse(responseText)
    };
  } catch (error) {
    return {
      success: false,
      error: error.toString()
    };
  }
}

function writeResponseToSheet(data) {
  const sheet = SpreadsheetApp.getActiveSheet();
  sheet.getRange('A1').setValue(JSON.stringify(data, null, 2));
}