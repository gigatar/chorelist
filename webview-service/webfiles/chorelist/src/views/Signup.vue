<template>
  <b-container fluid>
    <b-alert :variant="alert.variant" v-model="alert.show" dismissible>{{
      alert.text
    }}</b-alert>
    <b-row align-h="center">
      <b-col lg="3">
        <b-container fluid>
          <img src="@/assets/logo.png" />
          <form @submit.prevent="noop" v-if="!submitted">
            <b-form-group class="mb-1">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="organization"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="family-input"
                  v-model="signupRequest.family.name"
                  required
                  type="text"
                  placeholder="Enter Family Name"
                  autocomplete="nickname"
                  maxlength="128"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-form-group class="mb-1">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="person"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="display-input"
                  v-model="signupRequest.person.name"
                  required
                  type="text"
                  placeholder="Enter Display Name"
                  autocomplete="nickname"
                  maxlength="128"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-form-group class="mb-1">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="mail"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="email-input"
                  v-model="signupRequest.person.email"
                  required
                  type="email"
                  placeholder="Enter Email Address"
                  autocomplete="email"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-form-group class="mb-1">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="lock"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="password-input"
                  v-model="signupRequest.person.password"
                  required
                  type="password"
                  placeholder="Enter Password"
                  autocomplete="new-password"
                  minlength="8"
                ></b-form-input>
              </b-input-group>
              <Password
                v-model="signupRequest.person.password"
                strength-meter-only
              />
            </b-form-group>
            <b-button
              block
              type="submit"
              variant="primary"
              @click="signup"
              :disabled="signupEnabled"
              >Signup</b-button
            >
          </form>
          <b-container v-if="submitted">
            Your signup request has been submitted, please check your email to
            confirm signup.
          </b-container>
        </b-container>
      </b-col>
    </b-row>
  </b-container>
</template>
<script>
import Password from "vue-password-strength-meter";
export default {
  name: "Signup",
  components: { Password },
  data: () => ({
    submitted: false,
    signupRequest: {
      person: { email: null, name: null, password: null, type: "parent" },
      family: { name: null }
    },
    alert: { show: false, variant: "danger", text: "Test" }
  }),
  computed: {
    // Only enable the signup button when all fields have been entered.
    signupEnabled() {
      if (
        this.signupRequest.person.email &&
        this.signupRequest.person.password.length >= 8 &&
        this.signupRequest.person.name &&
        this.signupRequest.family.name
      ) {
        return false;
      }
      return true;
    }
  },
  methods: {
    noop() {
      /*Prevent form submit and bug*/
    },
    /**
     * Signup for a new account.
     *
     * If fields are validated, just return to prevent the unneccessary POST
     * to the backend.
     */
    signup() {
      if (
        !this.signupRequest.family.name ||
        !this.signupRequest.person.name ||
        !this.signupRequest.person.email ||
        !this.signupRequest.person.password.length >= 8
      ) {
        return;
      }
      this.$store
        .dispatch("signup", this.signupRequest)
        .then(() => {
          this.submitted = true;
        })
        .catch(err => {
          var errText = "";
          switch (err) {
            case 400:
              errText = "Invalid Input";
              break;
            case 409:
              errText = "Account or signup already exists.";
              break;
            default:
              errText = "An internal error happend, please try again.";
          }
          this.alert = { show: true, variant: "danger", text: errText };
        });
    }
  }
};
</script>
