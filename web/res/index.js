
let keyboardHolder = document.querySelector('#keyboard-holder');
let stickyHeader = document.querySelector('#sticky-header');

document.getElementById('opt-layout').addEventListener('change', (event) => {
    let v = event.currentTarget.value;
    let url = new URL(window.location.href);
    url.searchParams.set('layout', v);
    window.location.replace(url.toString());
});

document.getElementById('opt-captions').addEventListener('change', (event) => {
    let v = event.currentTarget.value;
    let kc = document.getElementsByClassName("txt");
    switch(v) {
        case "character":
            for(e of kc) {
                e.innerHTML = "..."
            }
            break;
        case "symbol_split":
            for(e of kc) {
                e.innerHTML = e.dataset.symbol
                    .replaceAll(/grave|acute|left|right|slash/g," $&") // TODO: smarter, more complete regex
                    .replaceAll("_"," ")
                    .replaceAll(" ","<br>")
            }
            break;
        default:
            for(e of kc) {
                e.innerHTML = e.dataset.symbol
            }
    }
});

document.getElementById("opt-captions-sym").selected = true;

