import {html} from 'htm/preact';
import PocketBase from '/dist/pocketbase.es.mjs'
import {ClickerStats} from "/components/clicker.js";
import {Pocket} from "./pocket.js";
import {Levels} from "./components/levels.js";
import {Auth} from "./components/auth.js";
import {useSignal} from '@preact/signals';
import {Updgrades} from "./components/upgrades.js";
import {useState} from 'preact/hooks';


export function App() {
    const count = useSignal(0);
    const pb = new PocketBase();
    const [user, setUser] = useState(null);

    function logout() {
        pb.authStore.clear();
        location.reload();
    }

    if (user) {
        return html`
            <${Pocket.Provider} value=${pb}>
                <div style="margin: 4px; display: flex; justify-content: end">
                    <a class="pure-button pure-button-primary" href="#" onclick=${logout}>Logout</a>
                </div>
                <div class="pure-g">
                    <${ClickerStats} count=${count}></ClickerStats>
                    <${Levels} count=${count}></Levels>
                    <${Updgrades} count=${count}></Updgrades>
                </div>
            </${Pocket.Provider}>`;
    }

    return html`
        <${Pocket.Provider} value=${pb}>
            <${Auth} setUser=${setUser}/>
        </${Pocket.Provider}>`;

}