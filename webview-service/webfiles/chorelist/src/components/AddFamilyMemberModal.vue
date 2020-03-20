<template>
  <b-modal
    id="addFamilyMember"
    title="Add Family Member"
    size="lg"
    no-close-on-backdrop
    no-close-on-esc
    :ok-disabled="formComplete"
    ok-title="Create Member"
    @ok="submit()"
    @cancel="resetForm()"
  >
    <b-alert :variant="alert.variant" v-model="alert.show" dismissible>{{
      alert.text
    }}</b-alert>

    <b-form-group class="mb-1" :disabled="createActive">
      <b-input-group size="md">
        <b-input-group-prepend is-text>
          <octicon name="person"></octicon>
        </b-input-group-prepend>
        <b-form-input
          id="newname-input"
          v-model="newUser.name"
          required
          type="text"
          placeholder="Enter Display Name"
          autocomplete="nickname"
          maxlength="128"
        ></b-form-input>
      </b-input-group>
    </b-form-group>
    <b-form-group class="mb-1" :disabled="createActive">
      <b-input-group size="md">
        <b-input-group-prepend is-text>
          <octicon name="mail" scale="1"></octicon>
        </b-input-group-prepend>
        <b-form-input
          id="newemail-input"
          v-model="newUser.email"
          required
          type="email"
          placeholder="Enter Email Address"
          autocomplete="email"
        ></b-form-input>
      </b-input-group>
    </b-form-group>
    <b-form-group class="mb-1" :disabled="createActive">
      <b-input-group size="md">
        <b-input-group-prepend is-text>
          <octicon name="organization" scale="1"></octicon>
        </b-input-group-prepend>
        <b-form-select
          v-model="newUser.type"
          :options="options"
        ></b-form-select>
      </b-input-group>
    </b-form-group>
    <b-form-group class="mb-1" :disabled="createActive">
      <b-input-group size="md">
        <b-input-group-prepend is-text>
          <octicon name="lock"></octicon>
        </b-input-group-prepend>
        <b-form-input
          id="newpassword-input"
          v-model="newUser.password"
          required
          type="password"
          placeholder="Enter New Password"
          autocomplete="new-password"
          minlength="8"
        ></b-form-input>
      </b-input-group>
    </b-form-group>
    <b-form-group class="mb-1" :disabled="createActive">
      <Password v-model="newUser.password" strength-meter-only />
    </b-form-group>
  </b-modal>
</template>
<script>
import Password from "vue-password-strength-meter";
import { validEmail } from "@/helpers/helpers.js";
export default {
  name: "AddFamilyMember",
  components: {
    Password
  },
  computed: {
    formComplete() {
      if (
        // this.newUser.name.length < 0 ||
        validEmail(this.newUser.email) === false ||
        (this.newUser.type !== "parent" && this.newUser.type !== "child") ||
        this.newUser.password.length < 8
      ) {
        return true;
      }

      return false;
    }
  },
  data: () => ({
    alert: { show: false, variant: "danger", text: "" },
    options: [
      { value: "parent", text: "Parent" },
      { value: "child", text: "Child" }
    ],
    newUser: { name: "", email: "", type: "child", password: "" },
    createActive: false
  }),
  methods: {
    createModal() {
      this.$bvModal.show("addFamilyMember");
    },
    submit() {
      this.$store
        .dispatch("addFamilyMember", this.newUser)
        .then(() => {
          this.$parent.getFamily(); // Reload the family table.
          this.$parent.alert = {
            show: true,
            variant: "success",
            text: 'Successfully added "' + this.newUser.name + '"'
          };
          this.resetForm();
        })
        .catch(error => {
          this.alert = {
            variant: "danger",
            text: ""
          };
          switch (error.status) {
            case 400:
              this.alert.text = "Invalid Input";
              break;
            case 401:
              this.alert.text = "Unauthorized";
              break;
            case 403:
              this.alert.text = "Ask your parents to do this";
              break;
            case 409:
              this.alert.text = "Email address already in use";
              break;
            default:
              this.alert.text = "An unknown error occured";
          }
          this.alert.show = true;

          this.$bvModal.show("addFamilyMember");
        });
    },
    resetForm() {
      this.newUser = { name: "", email: "", type: "child", password: "" };
    }
  }
};
</script>
