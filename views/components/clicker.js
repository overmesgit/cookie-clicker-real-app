import {html} from 'htm/preact';
import {useContext, useEffect} from 'preact/hooks';
import {Pocket} from "/pocket.js";
import {formatNumber} from "./format.js"


function Clicker(props) {
    const pb = useContext(Pocket);

    async function sendMessage() {
        pb.send('/click', {method: "POST", body: {"event": "click"}})
    }

    return html`
        <a id="click-button" onclick=${sendMessage} style="display: flex; justify-content: center">
            <img src="dist/universe-icon.svg" alt="icon" width="400px" height="400px"/>
        </a>`
}

function Counter({count}) {

    const pb = useContext(Pocket);
    useEffect(async () => {
        const userCounters = await pb.collection('counter').getList(1, 50, {
            filter: `user = "${pb.authStore.model.id}"`,
        });
        const initialCounter = userCounters.items[0];
        count.value = initialCounter.count;
        pb.collection('counter').subscribe(initialCounter.id, function (e) {
            count.value = e.record.count;
        });
    }, []);

    return html`
        <div id="cookies">${formatNumber(count)}</div>`
}

export function ClickerStats({count}) {
    return html`
        <div class="pure-u-1-3">
            <div class="column">
                <h2 style="text-align: center; margin-top: 0">Particles Clicker</h2>
                <${Clicker}/>
                <${Counter} count=${count}/>
            </div>
        </div>`
}