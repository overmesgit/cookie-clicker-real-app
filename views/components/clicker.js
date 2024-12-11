import {html} from 'htm/preact';
import {useContext, useEffect} from 'preact/hooks';
import {Pocket} from "/pocket.js";
import {formatNumber} from "./format.js"
import {useComputed, useSignal} from '@preact/signals';


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

function Leaderboard() {
    const userStats = useSignal(new Map());
    const userIDtoName = useSignal(new Map());
    const total = useComputed(() => {
        let value = userStats.value;
    });

    async function updateUserNames() {
        const allUsers = await pb.collection('users').getFullList();
        const res = new Map();
        allUsers.forEach((user) => {
            res.set(user.id, user.name)
        })
        userIDtoName.value = res;
    }

    useEffect(async () => {
        await updateUserNames();
    }, []);

    const pb = useContext(Pocket);
    useEffect(async () => {
        pb.collection('counter').subscribe('*', function (e) {
            if (e.record.user == pb.authStore.model.id) {
                return
            }
            if (!userIDtoName.value.has(e.record.user)) {
                updateUserNames();
                return;
            }
            const newMap = new Map(userStats.value);
            newMap.set(e.record.user, e.record.count);
            userStats.value = newMap;
        });

    }, []);

    return html`
        <div style="margin: 20px">
            <h2>Leaderboard</h2>
            <table class="pure-table pure-table-striped pure-table-hover">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Score</th>
                    </tr>
                </thead>
                <tbody>
                ${Array.from(userStats.value.keys().map((k) => {
                    return html`
                        <tr>
                            <td>${userIDtoName.value.get(k)}</td>
                            <td>${userStats.value.get(k)}</td>
                        </tr>`
                }))}
                </tbody>
            </table>
        </div>`
}

export function ClickerStats({count}) {
    return html`
        <div class="pure-u-1-3">
            <div class="column">
                <h2 style="text-align: center; margin-top: 0">Particles Clicker</h2>
                <${Clicker}/>
                <${Counter} count=${count}/>
                    <${Leaderboard}/>
            </div>
        </div>`
}