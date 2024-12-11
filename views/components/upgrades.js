import {Pocket} from "/pocket.js";
import {useContext, useEffect, useState} from 'preact/hooks';
import {html} from 'htm/preact';
import {useComputed, useSignal} from '@preact/signals';

import {formatNumber} from "./format.js"


function Updgrade({upgradeResp, count, allUpgrades}) {
    const pb = useContext(Pocket);
    let initialCount = upgradeResp.user_upgrade ? upgradeResp.user_upgrade.count : 0;
    const [upgradeCount, setUpgradeCount] = useState(initialCount);

    useEffect(() => {
        allUpgrades.value += upgradeCount * upgradeResp.upgrade.rate;
    }, [])

    const price = useSignal(upgradeResp.price);
    const isDisabled = useComputed(() => count.value < price.value);

    async function addLevel() {
        const res = await pb.send('/upgrade', {
            method: "POST",
            body: {"upgrade_id": upgradeResp.upgrade.id}
        })
        allUpgrades.value += upgradeResp.upgrade.rate;
        setUpgradeCount(res.user_upgrade.count)
        price.value = res.price
    }

    return html`
        <a class="level pure-g ${isDisabled.value ? 'disabled' : ''}" onclick=${addLevel}>
            <span class="pure-u-1-5">
                <img src="dist/upgrades/${upgradeResp.upgrade.tag}.svg" width="50px" height="50px"/>
            </span>
            <span class="pure-u-3-5">
                <div>${upgradeResp.upgrade.name}</div>
                <div>Rate: ${formatNumber(upgradeResp.upgrade.rate)}</div>
                <div>Price: ${formatNumber(price)}</div>
            </span>
            <span class="pure-u-1-5"
                  style="text-align: center; font-size: 36px">${upgradeCount}</span>
        </a>
    `
}

export function Updgrades({count}) {
    const pb = useContext(Pocket);
    const [upgrades, setUpgrades] = useState([]);
    const allUpgrades = useSignal(0);

    useEffect(async () => {
        const upgrades = await pb.send('/upgrade', {method: "GET"})
        setUpgrades(upgrades);
    }, []);


    return html`
        <div class="pure-u-1-3">
            <div class="column">
                <h2 style="text-align: center; margin-top: 0">Upgrades
                        (+${formatNumber(allUpgrades)} per
                    click)</h2>
                ${upgrades.map((upgradeResp) => Updgrade({upgradeResp, count, allUpgrades}))}
            </div>
        </div>
    `
}