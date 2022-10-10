package web

var indexTmplStr = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>i3keys</title>
	<style>
	{{.CSS}}
	</style>
</head>
<body>
	<div id="select-layout">
		<h2>Which layout do you use?</h2>
		<button class="layoutbtn" id="select-iso">ISO</button>
		<button class="layoutbtn" id="select-ansi">ANSI</button>
	</div>
	<div id="keyboard-holder">
	</div>
	<script>
	{{.JSData}}
	</script>
	<script>{{.JS}}</script>
</body>
</html>`

var indexTmplCSS = `:root {
	--key-size: calc(100vw / 24);
  }

  html, body {
	font-family: helveticaneue-light,helvetica neue light,helvetica neue,Helvetica,Arial,lucida grande,sans-serif;
	font-weight: 300;
  }

  #select-layout {
	  text-align: center;
  }

  .layoutbtn {
	  width: 45vw;
	  height: 5em;

	  background-color: #ff655c;
	  border: 1px solid #ff655c;
	  color: #fff;
	  font-size: 1.2em;
	  cursor: pointer;
  }

  .layoutbtn:hover {
	background-color: #e0453c;
	border-color: ##e0453c;
  }
  
  .keyboard {
	display: grid;
	grid-template-columns: repeat(1000, calc(var(--key-size) / 4 - 0.4vw));
	grid-template-rows: repeat(6, var(--key-size));
	grid-column-gap: 0.4vw;
	grid-row-gap: calc(var(--key-size) / 4);
	padding-bottom: 2em;
  }
  
  .key {
	border: 1px solid #000000;
	border-radius: 5px;
	background-color: #ffffff;
	font-size: 0.7vw;
	word-break: break-all;
	text-align: center;
	line-height: var(--key-size);
  }

  .key > .txt {
  pointer-events: none;
  }
  
  .emptySmall {
	grid-column: auto / span 2;
	border: 1px solid #ffffff;
  }
  
  .single {
	grid-column: auto / span 4;
  }
  
  .emptySingle {
	grid-column: auto / span 4;
	border: 1px solid #ffffff;
  }
  
  .modifier {
	grid-column: auto / span 5;
  }
  
  .onehalf {
	grid-column: auto / span 6;
  }
  
  .semidouble {
	grid-column: auto / span 7;
  }
  
  .double {
	grid-column: auto / span 8;
  }
  
  .semilarge {
	grid-column: auto / span 9;
  }
  
  .large {
	grid-column: auto / span 11;
  }
  
  .space {
	grid-column: auto / span 25;
  }
  
  .doubleY {
	grid-column: auto / span 4;
	grid-row: auto / span 2;
  }
  
  .enterUp {
	grid-column: auto / span 6;
  }
  
  .enterDown {
	border-top: none;
	border-radius: 0px;
	grid-column: auto / span 5;
	height: calc(var(--key-size) + (var(--key-size) / 4));
	margin-top: calc(-1 * (var(--key-size) / 4) - 1px);
  }
  
  .row-0 {
	grid-row-start: 1;
  }
  
  .row-1 {
	grid-row-start: 2;
  }
  
  .row-2 {
	grid-row-start: 3;
  }
  
  .row-3 {
	grid-row-start: 4;
  }
  
  .row-4 {
	grid-row-start: 5;
  }
  
  .row-5 {
	grid-row-start: 6;
  }
  
  .txt {
	display: inline-block;
	vertical-align: middle;
	line-height: normal;
  }

  .key-desc {
    font-weight: bold;
  }

  .key-modifier {
	background-color: #5cff87;
  }

  .key-used {
	background-color: #ff655c;
  }

  a {
  color: blue;
  }
`
var indexTmplJS = `/*
keyTypes = {
  'single': [1, 1], //one key
  'emptySingle': [1, 1], //empty one key
  'emptySmall': [0.5, 1], //gap between f-keys
  'double': [2, 1], //backspace
  'onehalf': [1.5, 1], //tab
  'doubleY': [1, 2], //numpad +, enter
  'semidouble': [1.75, 1], //caps
  'modifier': [1.25, 1], //ctrl
  'semilarge': [2.25, 1],
  'large': [2.75, 1], //right shift
  'space': [6.25, 1], //space
  'enterUp': [1.5, 1],  //enter
  'enterDown': [1.25, 1],  //enter
};
*/

ansi = [
    ["single", "emptySingle", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single"],
    ["single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "double", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "single"],
    ["onehalf", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "onehalf", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "doubleY"],
    ["semidouble", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "semilarge", "emptySmall", "emptySingle", "emptySingle", "emptySingle", "emptySmall", "single", "single", "single"],
    ["semilarge", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "large", "emptySmall", "emptySingle", "single", "emptySingle", "emptySmall", "single", "single", "single", "doubleY"],
    ["modifier", "modifier", "modifier", "space", "modifier", "modifier", "modifier", "modifier", "emptySmall", "single", "single", "single", "emptySmall", "double", "single"]
];

iso = [
    ["single", "emptySingle", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single"],
    ["single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "double", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "single"],
    ["onehalf", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "enterUp", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "doubleY"],
    ["semidouble", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "enterDown", "emptySmall", "emptySingle", "emptySingle", "emptySingle", "emptySmall", "single", "single", "single"],
    ["modifier", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "large", "emptySmall", "emptySingle", "single", "emptySingle", "emptySmall", "single", "single", "single", "doubleY"],
    ["modifier", "modifier", "modifier", "space", "modifier", "modifier", "modifier", "modifier", "emptySmall", "single", "single", "single", "emptySmall", "double", "single"]
];

let keyboardHolder = document.querySelector('#keyboard-holder');

function liLink(name, id) {
    let li = document.createElement('li');
    let a = document.createElement('a');
    a.href = "#"+id;
    a.title = name;
    a.innerText = name;
    li.appendChild(a);

    return li;
}

function generateKeyboard(layout, generated, modes) {
    let kbLayout;

    if (layout == "ANSI") {
        kbLayout = ansi;
    } else {
        kbLayout = iso;
    }
    
    let tosHeading = document.createElement('h1');
    tosHeading.innerHTML = "Table of contents";
    keyboardHolder.appendChild(tosHeading);
    let tosList = document.createElement('ul');
    let tosDefaultLi = liLink('Mode: Default', 'mode_head_default');
    tosDefaultUl = document.createElement('ul');
    
    for (let i = 0; generated !== null && i < generated.length; i++) {
        liEl = liLink(generated[i].Name, "keyboard_"+i);
        tosDefaultUl.appendChild(liEl);
    }
    tosList.appendChild(tosDefaultLi);
    tosList.appendChild(tosDefaultUl);
    keyboardHolder.appendChild(tosList);

    for (let i = 0; modes !== null && i < modes.length; i++) {
        let liEl = liLink("Mode: "+modes[i].Name, "mode_head_"+i);
        let ulEl = document.createElement('ul');

        for (let j = 0; j < modes[i].Keyboards.length; j++) {
            let liC = liLink(modes[i].Keyboards[j].Name, "mode_"+i+"_"+j);
            ulEl.appendChild(liC);
        }

        tosList.appendChild(liEl);
        tosList.appendChild(ulEl);
    }
    
    let headingEl = document.createElement('h2');
    headingEl.innerHTML = "Mode: Default";
    headingEl.id = "mode_head_default";
    keyboardHolder.appendChild(headingEl);

    for (let i = 0; generated !== null && i < generated.length; i++) {
        let newKeyboardGroup = generateKeyboardGroup(kbLayout, generated[i], "keyboard_"+i);
        keyboardHolder.appendChild(newKeyboardGroup);
    }


    for (let i = 0; modes !== null && i < modes.length; i++) {
        let headingEl = document.createElement('h2');
        headingEl.id = "mode_head_"+i;
        headingEl.innerHTML = "Mode: " + modes[i].Name;
        keyboardHolder.appendChild(headingEl);
        for (let j = 0; j < modes[i].Keyboards.length; j++) {
            let newKeyboardGroup = generateKeyboardGroup(kbLayout, modes[i].Keyboards[j], "mode_"+i+"_"+j);
            keyboardHolder.appendChild(newKeyboardGroup);
        }
    }
}

function generateKeyboardGroup(kbLayout, generated, headingID) {
    let kbWrapper = document.createElement('div');

    let headingEl = document.createElement('h3');
    headingEl.id = headingID;
    let headingContent = generated.Name;


    headingEl.innerHTML = headingContent;
    kbWrapper.appendChild(headingEl);

    let keyboardEl = document.createElement('div');
    keyboardEl.className = 'keyboard';

    let kbTextDescElement = document.createElement('div');
    let kbTextDescPElement = document.createElement('p');
    let kbTextDescResElement = document.createElement('span');
    kbTextDescResElement.className = "key-desc";
    kbTextDescPElement.innerHTML = "Hover over a key to see the command bound to the key here: ";
    kbTextDescPElement.appendChild(kbTextDescResElement);
    kbTextDescElement.appendChild(kbTextDescPElement);

    let enterHit = 0;

    for (let i = 0; i < kbLayout.length; i++) {
        let k = 0;
        for (let j = 0; j < kbLayout[i].length; j++) {
            keyEl = document.createElement('div');
            let gHit = 0;
            if (kbLayout[i][j] != "emptySingle" && kbLayout[i][j] != "emptySmall" && kbLayout[i][j] != "enterDown") {

                if (generated.Keys[i][k].Modifier) {
                    gHit = 1;
                }

                if (generated.Keys[i][k].InUse) {
                    gHit = 2;
                    keyEl.dataset.command = generated.Keys[i][k].Binding.command;

                    keyEl.addEventListener("mouseover", (e) => {
                        kbTextDescResElement.innerHTML = e.target.dataset.command;
                    });
                }

                if (kbLayout[i][j] == "enterUp") {
                    enterHit = gHit
                }

                keyEl.innerHTML = '<span class="txt">' + generated.Keys[i][k].Symbol + '</span>';

                k++;
            }

            if (kbLayout[i][j] == "enterDown") {
                gHit = enterHit;
            }

            let usedStatus = "";
            if (gHit == 1) {
                usedStatus = "key-modifier";
            } else if (gHit == 2) {
                usedStatus = 'key-used';
            }
            keyEl.className = "key " + kbLayout[i][j] + " row-" + i + " " + usedStatus;
            keyboardEl.appendChild(keyEl);
        }
    }

    kbWrapper.appendChild(keyboardEl);
    kbWrapper.appendChild(kbTextDescElement);

    return kbWrapper;
}

document.getElementById('select-iso').addEventListener('click', function () {
    generateKeyboard('ISO', generatedISO, generatedISOmodes);
    document.getElementById('select-layout').style.display = 'none';
});

document.getElementById('select-ansi').addEventListener('click', function () {
    generateKeyboard('ANSI', generatedANSI, generatedANSImodes);
    document.getElementById('select-layout').style.display = 'none';
});
`
