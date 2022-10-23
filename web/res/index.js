
let keyCaptions = {
    "space":        "Space",
    "exclam":       "!",
    "quotedbl":     "\"",
    "numbersign":   "#",
    "dollar":       "$",
    "percent":      "%",
    "ampersand":    "&amp;",
    "quoteright":   "'",
    "parenleft":    "(",
    "parenright":   ")",
    "asterisk":     "*",
    "plus":         "+",
    "comma":        ",",
    "minus":        "-",
    "period":       ".",
    "slash":        "/",
    "colon":        ":",
    "semicolon":    ";",
    "less":         "<",
    "equal":        "=",
    "greater":      ">",
    "question":     "?",
    "at":           "@",
    "bracketleft":  "[",
    "backslash":    "\\",
    "bracketright": "]",
    "asciicircum":  "^",
    "underscore":   "_",
    "quoteleft":    "&#96;",
    "braceleft":    "{",
    "bar":          "|",
    "braceright":   "}",
    "asciitilde":   "~",
    "agrave":       "à",
    "aacute":       "á",
    "acircumflex":  "â",
    "atilde":       "ã",
    "adiaeresis":   "ä",
    "aring":        "å",
    "ae":           "æ",
    "ccedilla":     "ç",
    "egrave":       "è",
    "eacute":       "é",
    "ecircumflex":  "ê",
    "ediaeresis":   "ë",
    "igrave":       "ì",
    "iacute":       "í",
    "icircumflex":  "î",
    "idiaeresis":   "ï",
    "eth":          "ð",
    "ntilde":       "ñ",
    "ograve":       "ò",
    "oacute":       "ó",
    "ocircumflex":  "ô",
    "otilde":       "õ",
    "odiaeresis":   "ö",
    "division":     "÷",
    "ooblique":     "ø",
    "ugrave":       "ù",
    "uacute":       "ú",
    "ucircumflex":  "û",
    "udiaeresis":   "ü",
    "yacute":       "ý",
    "thorn":        "þ",
    "ydiaeresis":   "ÿ",
    "KP_Multiply":  "KP *",
    "KP_Add":       "KP +",
    "KP_Subtract":  "KP -",
    "KP_Decimal":   "KP .",
    "KP_Divide":    "KP /",
    "KP_Equal":     "KP ="
}

function breakSpaces(s) { return s.replaceAll(" ","<br/>") }

function keyCaption(s) {
    let c = keyCaptions[s];
    return c ? c : splitSymbol(s);
}

function splitSymbol(s) {
    return s
        .replaceAll(/grave|acute|left|right|slash/g," $&") // TODO: smarter, more complete regex
        .replaceAll("_"," ")
        ;
}

let stickyHeader = document.getElementById('sticky-header');
let keyboardHolder = document.getElementById('keyboard-holder');
let keyInfoHolder = document.getElementById('key-info');
let keyElements = document.getElementsByClassName('key');

document.getElementById('opt-layout').addEventListener('change', (event) => {
    let v = event.currentTarget.value;
    let url = new URL(window.location.href);
    url.searchParams.set('layout', v);
    window.location.replace(url.toString());
});

document.getElementById('opt-captions').addEventListener('change', (event) => {
    let v = event.currentTarget.value;
    let kc = document.getElementsByClassName("txt");
    for(e of kc) {
        switch(v) {
            case "character":
                    e.innerHTML = breakSpaces(keyCaption(e.dataset.symbol));
                    e.classList.add("monospace")
                    break;
            case "symbol":
                    e.innerHTML = e.dataset.symbol;
                    e.classList.remove("monospace")
                    break;
            case "symbol_split":
                    e.innerHTML = breakSpaces(splitSymbol(e.dataset.symbol));
                    e.classList.remove("monospace")
                    break;
            case "symbol_code_dec":
                    e.innerHTML = e.dataset.symbolcode;
                    e.classList.add("monospace")
                    break;
            case "symbol_code_hex":
                    e.innerHTML = parseInt(e.dataset.symbolcode).toString(16).toUpperCase();
                    e.classList.add("monospace")
                    break;
            case "identifier":
                    e.innerHTML = e.dataset.identifier;
                    e.classList.add("monospace")
                    break;
            default:
                    e.innerHTML = "?"; // should not happen
        }
    }
});

// bubbling event for all the .key elements
keyboardHolder.addEventListener('mouseover', event => {
    for (el of keyElements) {
        if (el == event.target) {
            if (key = event.target.firstChild) { // key is actually .txt
                // TODO: somehow generate from range
                document.getElementById('key-command').innerHTML=key.dataset.bindcommand;
                document.getElementById('key-symbol').innerHTML=key.dataset.symbol;
                document.getElementById('key-symbolcode').innerHTML="0x"+parseInt(key.dataset.symbolcode).toString(16)+" ("+key.dataset.symbolcode+")";
                document.getElementById('key-identifier').innerHTML=key.dataset.identifier;
            } else {
                document.getElementById('key-command').innerHTML="";
                document.getElementById('key-symbol').innerHTML="";
                document.getElementById('key-symbolcode').innerHTML="";
                document.getElementById('key-identifier').innerHTML="";
            }
            break
        }
    }
});

document.getElementById("opt-captions").value = "symbol_split";
document.getElementById("opt-captions").dispatchEvent(new Event('change'));
