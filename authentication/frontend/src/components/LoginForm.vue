<template>
    <section id="loginView">
        <form id="loginForm" action="#" method="POST" @submit.prevent="login()">
            <p id="badlogin" v-bind:class="{ showing: loginStatus == 'failed' }">The username or password provided is incorrect.</p>
            <article class="input">
                <input v-model="username" name="username" autocomplete="username" id="username" placeholder="Username" required>
                <label for="username">Username</label>
            </article>
            <article class="input">
                <input v-model="password" name="password" autocomplete="current-password" id="password" type="password" placeholder="Password" required>
                <label for="password">Password</label>
            </article>
            <input type="submit" v-bind:value="loginStatus == 'unclicked' ? 'Login' : loginStatus == 'logging in' ? 'Logging in...' : loginStatus == 'success' ? 'Logged in successfully' : 'Login' ">
        </form>
        <form id="loginMfaForm" action="#" method="POST" @submit.prevent="mfa()">
            <p>It looks like your account has extra security features enabled. Please type in your TOTP token:</p>
            <article class="input">
                <input v-model="token" name="token" id="token" type="number" placeholder="TOTP Token" v-on:input="() => { verifyTotp() }" v-bind:disabled="tfaDisabled" required>
                <label for="token">TOTP Token</label>
            </article>
            <p id="verificationStatus" v-bind:class="{
                show: verificationStatus != 'unchecked'
            }">{{verificationStatus == 'failed' ? verificationFailReason : verificationStatus = 'verifying' ? "Verifying..." : ""}}</p>
        </form>
    </section>
</template>

<script>
import config from '../config'
export default {
  name: 'LoginForm',
  data: function() {
    return {
      loginFailed: false,
      loginFailReason: null,
      loginStatus: "unclicked",
      username: "",
      password: "",
      mfaSID: null,
      token: "",
      tfaDisabled: false,
      verificationStatus: 'unchecked',
      verificationFailReason: null
    }
  },
    methods: {
        verifyTotp: function() {
            if(this.token.toString().length != 6) {
                return;
            }

            this.verificationStatus = 'verifying'
            this.tfaDisabled = true;

            const data = JSON.stringify({
                "type":"totp",
                "username": this.username,
                "sid": this.mfaSID,
                "code": this.token
                });

            var xhr = new XMLHttpRequest();
            xhr.withCredentials = true;

            const _this = this;
            xhr.addEventListener("readystatechange", function() {
                if(this.readyState === 4) {
                        let resp = this.responseText

                        if (this.status != 200) {
                            switch(resp) {
                                case "invalid totp code":
                                    _this.verificationFailReason = "Invalid TOTP code provided."
                                    break;
                                case "totp not enabled":
                                    _this.verificationFailReason = "TOTP is not enabled on this account."
                                    break;
                                default:
                                    _this.verificationFailReason = `debug: ${resp}`

                            }
                            _this.verificationStatus = 'failed'
                            _this.tfaDisabled = false
                            return
                        }

                        _this.verificationStatus = 'success'

                        resp = JSON.parse(resp)

                        document.cookie = `sid=${resp['token']}; path=/`

                        alert("waiting")
                        window.location.href = "/account"
                }
            });

            xhr.open("POST", `${config.API_URL}/user/mfa/verify`);
            xhr.setRequestHeader("Content-Type", "application/json");

            xhr.send(data);
        },
        login: function() {
            this.loginStatus = "logging in"

            const data = JSON.stringify({"username":this.username,"password":this.password});

            var xhr = new XMLHttpRequest();
            xhr.withCredentials = true;

            const _this = this;
            xhr.addEventListener("readystatechange", function() {
                if(this.readyState === 4) {
                        let resp = this.responseText

                        if (this.status != 200) {
                            _this.loginStatus = 'failed'
                            _this.loginFailed = true
                            return
                        }

                        _this.loginStatus = 'success'

                        resp = JSON.parse(resp)

                        if(resp["mfaSID"] != undefined) {
                            const sid = resp["mfaSID"]
                            _this.mfaSID = sid;

                            document.querySelector("form#loginForm").classList.add("hide")
                            document.querySelector("form#loginMfaForm").classList.add("show")

                            return;
                        }

                        document.cookie = `sid=${resp['token']}; path=/`

                        alert("waiting")
                        window.location.href = "/account"
                }
            });

            xhr.open("POST", `${config.API_URL}/login`);
            xhr.setRequestHeader("Content-Type", "application/json");

            xhr.send(data);
      }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    @import 'styles.scss';

    #loginMfaForm {
        display: none;
    }
    #loginForm.hide {
        display: none;
    }

    p#verificationStatus {
        padding: 2px;
        margin: 0;
        color: rgb(17, 54, 219);
        display: none;
    }

    p#verificationStatus.verified {
        color: rgb(39, 206, 39);
    }

    p#verificationStatus.failed {
        color: rgb(240, 24, 24);
    }
    
    #loginMfaForm.show,
    p#verificationStatus.show {
        display: inherit;
    }
    
</style>
