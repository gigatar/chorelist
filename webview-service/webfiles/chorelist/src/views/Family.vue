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
            :disabled="modifyFamilyName || userType !== 'parent'"
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
                v-if="userType === 'parent'"
                >Update Name</b-button
              >
            </b-input-group>
          </b-form-group>
        </b-form></b-col
      ></b-row
    >
    <b-row
      ><b-col v-if="!loadingFamily">
        <b-table
          striped
          hover
          :items="familyMembers"
          :fields="familyFields"
          :busy="loadingFamily"
          sort-by="type"
          sort-desc
        >
          <template slot="table-busy" class="text-center text-primary my-2">
            <b-spinner type="grow" class="align-middle"></b-spinner>
            <strong>Loading Family...</strong>
          </template>
          <template v-slot:cell(actions)="item"
            ><b-button
              variant="outline-danger"
              @click="deleteFamilyMember(item.item.id, item.item.name)"
              v-if="userType === 'parent' && currentUserID !== item.item.id"
            >
              <octicon name="trashcan" scale="1"></octicon> </b-button
          ></template>
        </b-table>

        <b-button
          variant="primary"
          block
          @click="addFamilyMember()"
          v-if="userType === 'parent'"
          >Add Family Member</b-button
        >
      </b-col>
      <b-col v-if="loadingFamily">
        <div class="text-center text-primary my-2">
          <b-spinner type="grow" class="align-middle"></b-spinner>
          <strong>Loading Family...</strong>
        </div>
      </b-col></b-row
    ><b-row class="mt-5" v-if="!loadingFamily"
      ><b-col class="text-right">
        <b-button
          variant="danger"
          v-if="userType === 'parent'"
          @click="deleteFamily()"
          >Delete Family</b-button
        ></b-col
      ></b-row
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
  computed: {
    userType() {
      return this.$store.getters.getUserType;
    },
    currentUserID() {
      return this.$store.getters.getUserID;
    }
  },
  data: () => ({
    alert: { show: false, variant: "danger", text: "" },
    loadingFamily: true,
    modifyFamilyName: false,
    familyName: { name: null },
    familyFields: [
      { key: "name", sortable: true },
      { key: "type", sortable: true },
      { key: "lastLogin", sortable: true },
      { key: "Actions" }
    ],
    familyMembers: []
  }),
  beforeMount() {
    this.getFamily();
  },
  methods: {
    getFamily() {
      this.loadingFamily = true;
      this.familyName = { name: null };
      this.familyMembers = [];
      let refCount = 0;
      this.$store
        .dispatch("getFamily")
        .then(success => {
          this.familyName.name = success.data.name;

          for (let i = 0; i < success.data.person.length; i++) {
            refCount++;
            this.getFamilyMember(success.data.person[i])
              .then(response => {
                this.familyMembers.push({
                  id: response.id,
                  email: response.email,
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
                    this.alert.text = "Invalid Family or Family Member ID";
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
                refCount--;
                if (refCount === 0) {
                  this.loadingFamily = false;
                }
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
              this.alert.text = "An unknown error occured " + error;
          }
          this.alert.show = true;
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
            reject(error);
          });
      });
    },
    deleteFamily() {
      this.$bvModal
        .msgBoxConfirm(
          "THIS IS A DESTRUCTIVE AND IRREVERSABLE ACTION.  CONTINUING WILL DELETE ALL PEOPLE AND CHORES ASSOCIATED WITH THIS FAMILY.",
          {
            title: "Please Confirm",
            size: "lg",
            buttonSize: "sm",
            okVariant: "danger",
            okTitle: "YES",
            cancelTitle: "NO",
            footerClass: "p-2",
            hideHeaderClose: true,
            noCloseOnBackdrop: true,
            noCloseOnEsc: true
          }
        )
        .then(value => {
          // Even though we should only get bool bail if it's not explicitly true.
          if (value !== true) {
            return;
          }

          this.$store
            .dispatch("deleteFamily")
            .then(() => {
              this.$store.dispatch("logout");
            })
            .catch(error => {
              this.alert = {
                variant: "danger",
                text: ""
              };
              switch (error.status) {
                case 400:
                  this.alert.text = "Invalid Family or User ID";
                  break;
                case 401:
                  this.alert.text = "Unauthorized";
                  break;
                case 403:
                  this.alert.text = "HA! Nice try kid.";
                  break;
                default:
                  this.alert.text = "An unknown error occured";
              }
              this.alert.show = true;
            });
        });
    },
    deleteFamilyMember(userID, userName) {
      this.$bvModal
        .msgBoxConfirm(
          "Are you sure you want to delete " +
            userName +
            "?  This action is irreversable and will remove them as the assignee of all chores.",
          {
            title: "Please Confirm",
            size: "lg",
            buttonSize: "sm",
            okVariant: "danger",
            okTitle: "YES",
            cancelTitle: "NO",
            footerClass: "p-2",
            hideHeaderClose: true,
            noCloseOnBackdrop: true,
            noCloseOnEsc: true
          }
        )
        .then(value => {
          // Even though we should only get bool bail if it's not explicitly true.
          if (value !== true) {
            return;
          }

          this.$store
            .dispatch("deleteFamilyMember", userID)
            .then(() => {
              this.getFamily();
              this.alert = {
                variant: "success",
                show: true,
                text: "Successfully deleted " + userName
              };
            })
            .catch(error => {
              this.alert = {
                variant: "danger",
                text: ""
              };
              switch (error.status) {
                case 400:
                  this.alert.text = "Invalid Family or User ID";
                  break;
                case 401:
                  this.alert.text = "Unauthorized";
                  break;
                case 403:
                  this.alert.text = "HA! Nice try kid.";
                  break;
                default:
                  this.alert.text = "An unknown error occured";
              }
              this.alert.show = true;
            });
        });
    }
  }
};
</script>
