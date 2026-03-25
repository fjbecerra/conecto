function showSidebar() {
  const htmlServ = HtmlService.createTemplateFromFile('index');
  const html = htmlServ.evaluate();
  const ui = SpreadsheetApp.getUi();
  ui.showSidebar(html);

}

