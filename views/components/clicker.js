import {html} from 'htm/preact';
import {useContext, useEffect} from 'preact/hooks';
import {Pocket} from "/pocket.js";
import {formatNumber} from "./format.js"
import {useComputed, useSignal} from '@preact/signals';


function Clicker() {
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
    return html`
        <div id="cookies">${formatNumber(count)}</div>`
}

function Leaderboard({count}) {
    const pb = useContext(Pocket);
    const userStats = useSignal(new Map());
    const sortedStats = useComputed(() => {
        const data = Array.from(userStats.value.keys().map((k) => {
            return [userIDtoName.value.get(k), userStats.value.get(k), k]
        }))
        data.sort((a, b) => b[1] - a[1])
        return data;
    })
    const userIDtoName = useSignal(new Map());

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
        const allCounter = await pb.collection('counter').getFullList();
        const res = new Map();
        allCounter.forEach((conter) => {
            res.set(conter.user, conter.count)
        })
        count.value = res.get(pb.authStore.model.id);
        userStats.value = res;
    }, []);


    useEffect(async () => {
        pb.collection('counter').subscribe('*', function (e) {
            if (!userIDtoName.value.has(e.record.user)) {
                updateUserNames();
                return;
            }
            const newMap = new Map(userStats.value);
            newMap.set(e.record.user, e.record.count);
            userStats.value = newMap;
            count.value = newMap.get(pb.authStore.model.id);
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
                ${sortedStats.value.map(([name, value, userID]) => {
                    const font = userID === pb.authStore.model.id ? 'bold' : 'normal';
                    return html`
                        <tr style="font-weight: ${font}">
                            <td>${name}</td>
                            <td>${formatNumber(value)}</td>
                        </tr>`
                })}
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
                    <${Leaderboard} count=${count}/>
            </div>
        </div>`
}