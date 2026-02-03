const { app, BrowserWindow } = require('electron');
const hello = require('./output/hello.node');

console.log('Electron:', process.versions.electron);
console.log('Node:', process.versions.node);
console.log('ABI:', process.versions.modules);
console.log('');
console.log('Testing hello.node:');
console.log('Hello1:', hello.Hello1('test', '10'));
console.log('Hello2:', hello.Hello2('test', '10'));
console.log('ReturnString:', hello.ReturnString('test', '10'));
console.log('ReturnInt:', hello.ReturnInt('test', '10'));
console.log('ReturnFloat:', hello.ReturnFloat('test', '10'));
console.log('ReturnBool:', hello.ReturnBool('test', '10'));
console.log('ReturnObject:', JSON.stringify(hello.ReturnObject('test', '10')));
console.log('ReturnNestedObject:', JSON.stringify(hello.ReturnNestedObject('test', '10')));

app.whenReady().then(() => {
    const win = new BrowserWindow({ width: 800, height: 600 });
    win.loadFile('index.html');
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});
