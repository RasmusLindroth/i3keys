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
	background-color: #ffffff;
	font-size: 0.7vw;
	word-break: break-all;
	text-align: center;
	line-height: var(--key-size);
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

  .key-modifier {
	background-color: #5cff87;
  }

  .key-used {
	background-color: #ff655c;
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


function generateKeyboard(generated) {

    let kbLayout;

    if (generated.type == "ANSI") {
        kbLayout = ansi;
    } else {
        kbLayout = iso;
    }

    for (let i = 0; i < groups.length; i++) {
        let newKeyboardGroup = generateKeyboardGroup(kbLayout, generated, groups[i]);
        keyboardHolder.appendChild(newKeyboardGroup);
    }
}

function generateKeyboardGroup(kbLayout, generated, group) {
    let kbWrapper = document.createElement('div');

    let headingEl = document.createElement('h3');
    let headingContent = 'No modifiers';

    if (group.modifiers != null && group.modifiers.length > 0) {
        headingContent = group.modifiers.join('+');
    }

    headingEl.innerHTML = headingContent;
    kbWrapper.appendChild(headingEl);

    let keyboardEl = document.createElement('div');
    keyboardEl.className = 'keyboard';

    modifierListKeys = Object.keys(modifierList);

    let enterHit = 0;

    for (let i = 0; i < kbLayout.length; i++) {
        let k = 0;
        for (let j = 0; j < kbLayout[i].length; j++) {
            keyEl = document.createElement('div');
            let gHit = 0;
            if (kbLayout[i][j] != "emptySingle" && kbLayout[i][j] != "emptySmall" && kbLayout[i][j] != "enterDown") {

                for (let mi = 0; group.modifiers != null && mi < group.modifiers.length; mi++) {
                    if (!modifierListKeys.includes(group.modifiers[mi])) {
                        continue;
                    }
                    for (let mj = 0; mj < modifierListKeys.length; mj++) {
                        let mk = modifierListKeys[mj];
                        if (mk != group.modifiers[mi]) {
                            continue;
                        }

                        for (let mx = 0; mx < modifierList[mk].length; mx++) {
                            if (generated.content[i][k] == modifierList[mk][mx]) {
                                gHit = 1;
                            }
                        }
                    }
                }

                for (let gi = 0; gi < group.bindings.length; gi++) {
                    if (group.bindings[gi].key == generated.content[i][k]) {
                        gHit = 2;
                        break;
                    }
                }

                if (kbLayout[i][j] == "enterUp") {
                    enterHit = gHit
                }

                if (typeof generated.keys[i][k] !== 'undefiend' && generated.keys[i][k][0] == "A" && !blacklist.includes(generated.codes[i][k])) {
                    keyEl.innerHTML = "<span class=\"txt\">&#" + generated.codes[i][k] + ";</span>";
                } else {
                    keyEl.innerHTML = '<span class="txt">' + generated.content[i][k] + '</span>';
                }
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

    return kbWrapper;
}

document.getElementById('select-iso').addEventListener('click', function () {
    generateKeyboard(generatedISO);
    document.getElementById('select-layout').style.display = 'none';
});

document.getElementById('select-ansi').addEventListener('click', function () {
    generateKeyboard(generatedANSI);
    document.getElementById('select-layout').style.display = 'none';
});
`
