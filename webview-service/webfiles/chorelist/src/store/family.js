import axios from "axios";

export default {
  state: {},
  getters: {},
  mutations: {},
  actions: {
    changeFamilyName(context, newName) {
      return new Promise((resolve, reject) => {
        axios
          .patch("rest/v1/families/name", newName, {
            headers: {
              Authorization: "Bearer " + context.getters.getAuthToken
            }
          })
          .then(({ status }) => {
            resolve(status);
          })
          .catch(function(error) {
            reject(error.response);
          });
      });
    },
    getFamily(context){
        return new Promise((resolve,reject)=>{
            axios.get("/rest/v1/families",{
                headers: {
                    Authorization: "Bearer " + context.getters.getAuthToken
                }
            })
            .then(({data, status})=>{
                resolve({status,data});
            })
            .catch(function(error) {
                reject(error.response);
            })
        });
    }
  }
};
