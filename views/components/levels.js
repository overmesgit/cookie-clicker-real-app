import {Pocket} from "/pocket.js";
import {useContext, useEffect, useState} from 'preact/hooks';
import {html} from 'htm/preact';
import {useComputed, useSignal} from '@preact/signals';
import {formatNumber} from "./format.js"

function calculatePrice(levelCost, count) {
    return (levelCost * 1.3 ** count).toFixed(0);
}

function Level({levelResp, count, allLevels}) {
    function updateAllLevels(levelsSignal, value) {
        let newMap = new Map(levelsSignal.value);
        newMap.set(levelResp.level.id, value);
        levelsSignal.value = newMap;
    }

    const pb = useContext(Pocket);
    const [userLevelID, setUserLevelID] = useState(levelResp.user_level ? levelResp.user_level.id : null);
    const [userLevelCount, setUserLevelCount] = useState(levelResp.user_level ? levelResp.user_level.count : 0);
    useEffect(() => {
        let initialRate = userLevelCount * levelResp.level.rate;
        updateAllLevels(allLevels, initialRate);
    }, [])
    const price = useSignal(calculatePrice(levelResp.level.cost, userLevelCount));
    const isDisabled = useComputed(() => count.value < price);

    useEffect(async () => {
        if (userLevelID) {
            await pb.collection('user_levels').subscribe(userLevelID, function (e) {
                setUserLevelCount(e.record.count);
                updateAllLevels(allLevels, levelResp.level.rate * e.record.count);
                price.value = calculatePrice(levelResp.level.cost, e.record.count)
            });
        }
    }, [userLevelID]);

    async function addLevel() {
        const res = await pb.send('/level', {
            method: "POST",
            body: {"level_id": levelResp.level.id}
        })

        if (!userLevelID) {
            updateAllLevels(allLevels, levelResp.level.rate * res.user_level.count);
            setUserLevelID(res.user_level.id);
            setUserLevelCount(res.user_level.count);
        }
    }

    function getIconName(name) {
        return name.replace(" ", "-").toLowerCase();
    }

    return html`
        <a class="level pure-g ${isDisabled.value ? 'disabled' : ''}" onclick=${addLevel}>
            <span class="pure-u-1-5">
                <img src="dist/levels/${getIconName(levelResp.level.name)}.svg" width="50px"
                     height="50px"/>
            </span>
            <span class="pure-u-3-5">
                <div>${levelResp.level.name}</div>
                <div>Rate: ${formatNumber(levelResp.level.rate)} per second (30 sec)</div>
                <div>Price: ${formatNumber(price)}</div>
            </span>
            <span class="pure-u-1-5"
                  style="text-align: center; font-size: 36px">${userLevelCount}</span>
        </a>
    `
}

export function Levels({count}) {
    const pb = useContext(Pocket);
    const [levels, setLevels] = useState([]);
    const allLevels = useSignal(new Map());
    const total = useComputed(() => allLevels.value.values().reduce((a, b) => a + b, 0));

    useEffect(async () => {
        const levels = await pb.send('/level', {method: "GET"})
        setLevels(levels);
    }, []);


    return html`
        <div class="pure-u-1-3">
            <div class="column">
                <h2 style="text-align: center; margin-top: 0">
                    Levels (+${formatNumber(total)} per sec)</h2>
                ${levels.map((levelResp) => Level({levelResp, count, allLevels}))}
            </div>
        </div>
    `
}