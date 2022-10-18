/*
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
let stickyHeader = document.querySelector('#sticky-header');

function a_name() {
    return [...arguments].join("_")
}

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

    // #sticky-header begin ---

    let tosHeading = document.createElement('h1');
    tosHeading.innerHTML = "Table of contents";
    stickyHeader.appendChild(tosHeading);

    let tosList = document.createElement('ul');
    let tosDefaultLi = liLink("Mode: "+"Default", a_name("keyboard",0));
    tosDefaultUl = document.createElement('ul');

    for (let i = 0; generated !== null && i < generated.length; i++) {
        liEl = liLink(generated[i].Name, a_name("keyboard",i));
        tosDefaultUl.appendChild(liEl);
    }
    tosList.appendChild(tosDefaultLi);
    tosDefaultLi.appendChild(tosDefaultUl);
    stickyHeader.appendChild(tosList);

    for (let i = 0; modes !== null && i < modes.length; i++) {
        let liEl = liLink("Mode: "+modes[i].Name, a_name("mode",i,0));
        let ulEl = document.createElement('ul');

        for (let j = 0; j < modes[i].Keyboards.length; j++) {
            let liC = liLink(modes[i].Keyboards[j].Name, a_name("mode",i,j));
            ulEl.appendChild(liC);
        }

        tosList.appendChild(liEl);
        liEl.appendChild(ulEl);
    }

    // options panel

    let optSplitLabel = document.createElement('label');
        optSplitLabel.for = "opt-split";
        optSplitLabel.innerText = "Split"

    let optSplitInput = document.createElement('input');
        optSplitInput.type = "checkbox";
        optSplitInput.name = "opt-split";
        optSplitInput.id = "opt-split";
        optSplitInput.checked = true;

        optSplitInput.addEventListener('change', (event) => {
            let split = event.currentTarget.checked;
            let kc = document.getElementsByClassName("txt");
            console.log('opt-split','change',event,split,kc.length);
            for (let i = 0; i < kc.length; i++) {
                if (split) {
                    kc[i].innerHTML = kc[i].innerHTML.replaceAll("_","<br>");
                } else {
                    kc[i].innerHTML = kc[i].innerHTML.replaceAll("<br>","_");
                }
            }
        });

    let optPanel = document.createElement('div');    
        optPanel.id = "options-panel";

    optPanel.appendChild(optSplitLabel);
    optPanel.appendChild(optSplitInput);
    stickyHeader.appendChild(optPanel);
    
    /*
    <div id="options-panel">
    <label for="opt-split">Split</label><br></br>
    <input type="checkbox" id="opt-split" name="opt-split" value="1"></input>
    </div>
    */

    // #sticky-header end ---

    /*
    let headingEl = document.createElement('h2');
    headingEl.innerHTML = /*"Mode: Default";
    headingEl.id = a_name("mode","default");
    keyboardHolder.appendChild(headingEl);
    */

    for (let i = 0; generated !== null && i < generated.length; i++) {
        let newKeyboardGroup = generateKeyboardGroup(kbLayout, generated[i], "Default", a_name("keyboard",i));
        keyboardHolder.appendChild(newKeyboardGroup);
    }

    for (let i = 0; modes !== null && i < modes.length; i++) {
        for (let j = 0; j < modes[i].Keyboards.length; j++) {
            let newKeyboardGroup = generateKeyboardGroup(kbLayout, modes[i].Keyboards[j], modes[i].Name, a_name("mode",i,j));
            keyboardHolder.appendChild(newKeyboardGroup);
        }
    }
}

function generateKeyboardGroup(kbLayout, generated, modeName, headingID) {
    let kbWrapper = document.createElement('div');

    let headingWrapper = document.createElement('div');
    headingWrapper.className = "keyboard-heading";

    let headingEl2 = document.createElement('h2');
    /*headingEl2.id = "mode_head_"+i;*/
    headingEl2.innerHTML = "Mode: " + modeName;
    headingWrapper.appendChild(headingEl2);

    let headingEl3 = document.createElement('h3');
    headingEl3.id = headingID;
    let headingContent = generated.Name;
    headingEl3.innerHTML = headingContent;
    headingWrapper.appendChild(headingEl3);

    kbWrapper.appendChild(headingWrapper);

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

                keyEl.innerHTML = '<span class="txt">' + generated.Keys[i][k].Symbol
                    .replaceAll(/grave|acute|left|right|slash/g," $&") // TODO: smarter, more complete regex
                    .replaceAll("_"," ")
                    .replaceAll(" ","<br>")
                + '</span>';

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