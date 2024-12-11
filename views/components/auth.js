import {html} from 'htm/preact';
import {useContext, useState} from 'preact/hooks';
import {Pocket} from "/pocket.js";


export function Auth({setUser}) {
    const [tab, setTab] = useState('login');

    const content = tab === 'login' ?
        html`
            <${Login} setUser=${setUser}/>` :
        html`
            <${SignUp} setUser=${setUser}/>`;

    const classes = {
        'login': tab === 'login' ? 'pure-button-active' : '',
        'signup': tab === 'signup' ? 'pure-button-active' : '',
    }
    return html`
        <div style="margin-top: 300px; display: flex; justify-content: center">
            <div>
                <div style="display: flex; justify-content: center; margin: 10px">
                    <div class="pure-button-group" role="group" aria-label="...">
                        <button class="pure-button ${classes['login']}"
                                onclick=${() => setTab('login')}>
                            Login
                        </button>
                        <button class="pure-button ${classes['signup']}"
                                onclick=${() => setTab('signup')}>
                            Sign
                            Up
                        </button>
                    </div>
                </div>
                ${content}
            </div>
        </div>
    `
}

function Login({setUser}) {
    const pb = useContext(Pocket);
    const [error, setError] = useState(null);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        let user = formData.get('email');
        try {
            await pb.collection("users").authWithPassword(user, formData.get('password'));
            setUser(user);
        } catch (error) {
            setError(error.message);
        }

    };

    const errorContent = error ? html`<p style="color: darkred; margin-left: 60px">
        ${error}</p>` : null;
    return html`
        <div>
            ${errorContent}
            <form class="pure-form pure-form-aligned" onSubmit=${handleSubmit}>
                <fieldset>
                    <div class="pure-control-group">
                        <label for="aligned-email">Email Address</label>
                        <input name="email" type="email" id="aligned-email"
                               placeholder="Email Address"/>
                    </div>
                    <div class="pure-control-group">
                        <label for="password">Password</label>
                        <input name="password" type="password" id="password"
                               placeholder="Password"/>
                    </div>
                    <div class="pure-controls">
                        <button type="submit" class="pure-button pure-button-primary">Submit
                        </button>
                    </div>
                </fieldset>
            </form>
        </div>`
}

function SignUp({setUser}) {
    const pb = useContext(Pocket);
    const [error, setError] = useState(null);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        setError(null);

        let user = formData.get('email');
        try {
            await pb.send('/signup', {
                method: "POST", body: {
                    "email": user,
                    "name": formData.get('name'),
                    "password": formData.get('password'),
                    "passwordConfirm": formData.get('passwordConfirm'),
                }
            });
            await pb.collection("users").authWithPassword(user, formData.get('password'));
            setUser(user);
        } catch (error) {
            setError(JSON.stringify(error.response));
        }
    };

    const errorContent = error ? html`<p style="color: darkred; margin-left: 60px">
        ${error}</p>` : null;

    return html`
        <div>
            ${errorContent}
            <form class="pure-form pure-form-aligned" onSubmit=${handleSubmit}>
                <fieldset>
                    <div class="pure-control-group">
                        <label for="name">Username</label>
                        <input name="name" type="text" id="name"
                               placeholder="Username"/>
                    </div>
                    <div class="pure-control-group">
                        <label for="aligned-email">Email Address</label>
                        <input name="email" type="email" id="aligned-email"
                               placeholder="Email Address"/>
                    </div>
                    <div class="pure-control-group">
                        <label for="password">Password</label>
                        <input name="password" type="password" id="password"
                               placeholder="Password"/>
                    </div>
                    <div class="pure-control-group">
                        <label for="confirm-password">Confirm Password</label>
                        <input name="passwordConfirm" type="password" id="confirm-password"
                               placeholder="Confirm Password"/>
                    </div>
                    <div class="pure-controls">
                        <button type="submit" class="pure-button pure-button-primary">Submit
                        </button>
                    </div>
                </fieldset>
            </form>
        </div>
    `
}