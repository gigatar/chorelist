<template>
  <b-modal
    id="userModal"
    title="User Profile"
    size="lg"
    no-close-on-backdrop
    no-close-on-esc
  >
    <b-container fluid>
      <b-alert :variant="alert.variant" v-model="alert.show" dismissible>{{
        alert.text
      }}</b-alert>
      <b-row>
        <b-col>
          <b-form>
            <b-form-group
              class="mb-3"
              label-size="lg"
              label="Change user name"
              :disabled="modifyUser"
            >
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="person"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="display-input"
                  v-model="person.name"
                  required
                  type="text"
                  placeholder="Enter Display Name"
                  autocomplete="nickname"
                  maxlength="128"
                ></b-form-input>
                <b-button
                  type="button"
                  class="ml-2"
                  variant="primary"
                  @click="updateName"
                  >Update Name</b-button
                >
              </b-input-group>
            </b-form-group>
          </b-form>
        </b-col>
      </b-row>
      <b-row
        ><b-col><hr /></b-col
      ></b-row>
      <b-row>
        <b-col>
          <b-form>
            <b-form-group
              class="mb-1 mt-2"
              label-size="lg"
              label="Change password"
              :disabled="changePassword"
            >
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="lock"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="oldpassword-input"
                  v-model="passwordUpdate.oldPassword"
                  required
                  type="password"
                  placeholder="Enter Current Password"
                  autocomplete="current-password"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-form-group class="mb-1" :disabled="changePassword">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="lock"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="newpassword-input"
                  v-model="passwordUpdate.password"
                  required
                  type="password"
                  placeholder="Enter New Password"
                  autocomplete="new-password"
                  minlength="8"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-form-group class="mb-1" :disabled="changePassword">
              <Password v-model="passwordUpdate.password" strength-meter-only />
              <b-button
                :disabled="
                  passwordUpdate.password.length < 8 ||
                    !passwordUpdate.oldPassword
                "
                type="button"
                class="ml-2"
                variant="primary"
                @click="updatePassword"
                block
                >Change Password</b-button
              >
            </b-form-group>
          </b-form></b-col
        ></b-row
      >
    </b-container>
  </b-modal>
</template>

<script>
import Password from "vue-password-strength-meter";
export default {
  name: "UserModal",
  components: {
    Password
  },
  data: () => ({
    person: { name: null },
    passwordUpdate: { password: "", oldPassword: "" },
    modifyUser: false,
    changePassword: false,
    alert: { show: false, variant: "danger", text: "" }
  }),
  methods: {
    createModal() {
      this.person.name = this.$store.getters.getUserName;
      this.$bvModal.show("userModal");
    },
    updateName() {
      // Disable form
      this.modifyUser = true;

      this.$store
        .dispatch("changeUserName", this.person)
        .then(() => {
          this.alert = {
            show: true,
            variant: "success",
            text: "Successfully changed name"
          };
        })
        .catch(error => {
          this.alert = {
            variant: "danger",
            text: ""
          };
          switch (error.status) {
            case 400:
              this.alert.text = "Invalid name";
              break;
            default:
              this.alert.text = "An unknown error occured.";
          }
          this.alert.show = true;
        })
        .finally(() => {
          this.modifyUser = false;
        });
    },
    updatePassword() {
      // Disable form
      this.changePassword = true;

      this.$store
        .dispatch("changePassword", this.passwordUpdate)
        .then(() => {
          this.alert = {
            show: true,
            variant: "success",
            text: "Successfully changed password"
          };
        })
        .catch(error => {
          this.alert = {
            variant: "danger",
            text: ""
          };
          switch (error.status) {
            case 400:
              this.alert.text = "New password doesn't meet requirements";
              break;
            case 403:
              this.alert.text = "Invalid old password";
              break;
            default:
              this.alert.text = "An unknown error occured.";
          }
          this.alert.show = true;
        })
        .finally(() => {
          this.passwordUpdate.password = "";
          this.passwordUpdate.oldPassword = "";
          this.changePassword = false;
        });
    }
  }
};
</script>
