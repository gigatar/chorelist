<template>
  <b-container>
    <Navigation />
    <b-alert :variant="alert.variant" v-model="alert.show" dismissible>{{
      alert.text
    }}</b-alert>

    <b-row>
      <b-col xs="12"
        ><b-form>
          <b-form-group
            class="mb-3"
            label-size="lg"
            label="Family Name"
            :disabled="modifyFamilyName"
          >
            <b-input-group size="md">
              <b-input-group-prepend is-text>
                <octicon name="organization"></octicon>
              </b-input-group-prepend>
              <b-form-input
                id="display-input"
                v-model="familyName.name"
                required
                type="text"
                placeholder="Enter Family Name"
                autocomplete="nickname"
                maxlength="128"
              ></b-form-input>
              <b-button
                type="button"
                class="ml-2"
                variant="primary"
                @click="updateFamilyName"
                >Update Name</b-button
              >
            </b-input-group>
          </b-form-group>
        </b-form></b-col
      ></b-row
    >
    <b-row
      ><b-col>
        <b-table
          striped
          hover
          :items="familyMembers"
          :fields="familyFields"
          :busy="loadingFamily"
        >
          <div slot="table-busy" class="text-center text-primary my-2">
            <b-spinner type="grow" class="align-middle"></b-spinner>
            <strong>Loading Family...</strong>
          </div></b-table
        >
        <b-button variant="primary" block @click="addFamilyMember()"
          >Add Family Member</b-button
        >
      </b-col></b-row
    >
    <AddFamilyMember ref="child"></AddFamilyMember>
  </b-container>
</template>

<script>
import Navigation from "@/components/Navigation";
import AddFamilyMember from "@/components/AddFamilyMemberModal";
import { unixtTimeToLocal } from "@/helpers/helpers.js";
export default {
  name: "Family",
  components: {
    Navigation,
    AddFamilyMember
  },
  data: () => ({
    alert: { show: false, variant: "danger", text: "" },
    loadingFamily: true,
    modifyFamilyName: false,
    familyName: { name: null },
    familyFields: [
      { key: "name", sortable: true },
      { key: "type", sortable: true },
      { key: "lastLogin", sortable: true }
    ],
    familyMembers: []
  }),
  beforeMount() {
    this.getFamily();
  },
  methods: {
    getFamily() {
      this.familyName = { name: null };
      this.familyMembers = [];
      this.$store
        .dispatch("getFamily")
        .then(success => {
          this.familyName.name = success.data.name;

          for (let i = 0; i < success.data.person.length; i++) {
            this.getFamilyMember(success.data.person[i])
              .then(response => {
                this.familyMembers.push({
                  name: response.name,
                  type: response.type,
                  lastLogin: unixtTimeToLocal(response.lastLogin) || "Never"
                });
              })
              .catch(error => {
                this.alert = {
                  variant: "danger",
                  text: ""
                };
                switch (error.status) {
                  case 400:
                    this.alert.text = "Invalid Family ID";
                    break;
                  case 401:
                    this.alert.text = "Unauthorized";
                    break;
                  default:
                    this.alert.text = "An unknown error occured";
                }
                this.alert.show = true;
              });
          }
        })
        .catch(error => {
          this.alert = {
            variant: "danger",
            text: ""
          };
          switch (error.status) {
            case 400:
              this.alert.text = "Invalid Family ID";
              break;
            case 401:
              this.alert.text = "Unauthorized";
              break;
            default:
              this.alert.text = "An unknown error occured";
          }
          this.alert.show = true;
        })
        .finally(() => {
          this.loadingFamily = false;
        });
    },
    addFamilyMember() {
      this.$refs.child.createModal();
    },
    updateFamilyName() {
      this.modifyFamilyName = true;
      this.$store
        .dispatch("changeFamilyName", this.familyName)
        .then(() => {
          this.alert = {
            variant: "success",
            text: "Successfully change family name",
            show: true
          };
        })
        .catch(error => {
          switch (error.status) {
            case 400:
              this.alert.text = "Invalid Family Name";
              break;
            case 401:
              this.alert.text = "Unauthorized";
              break;
            case 403:
              this.alert.text =
                "Do your parents know you're trying to change the name?";
              break;
            default:
              this.alert.text = "An unknown error occured";
          }
          this.alert.show = true;
        })
        .finally(() => (this.modifyFamilyName = false));
    },
    getFamilyMember(userID) {
      return new Promise((resolve, reject) => {
        this.$store
          .dispatch("getUser", userID)
          .then(success => {
            resolve(success);
          })
          .catch(error => {
            this.alert = {
              variant: "danger",
              text: ""
            };
            switch (error.status) {
              case 400:
                this.alert.text = "Invalid User ID";
                break;
              case 401:
                this.alert.text = "Unauthorized";
                break;
              default:
                this.alert.text = "An unknown error occured";
            }
            this.alert.show = true;
            reject();
          });
      });
    }
  }
};
</script>
