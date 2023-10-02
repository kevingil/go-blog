// Generate mock username
const randomNumber = Math.floor(Math.random() * (999999 - 100000 + 1)) + 100000;
let userName = "user" + randomNumber;

// User presses publish
document.getElementById('postButton').onclick = function() {
    postInput = document.getElementById("post-input");
    postContent = postInput.value;
    handleUserPost(userName, postContent)
    postInput.value = '';
  };


function handleUserPost(user, content){
  var postID = generatePostID();
  addToTimeline(user, content, postID);
  useModerator(content, postID);
}

var generatePostID = function(){
  const randomNumber = Math.floor(Math.random() * (999999 - 100000 + 1)) + 100000;
  let postID = "postID" + randomNumber;
  return postID;
}

function showLoadingAnimation() {

  }
  
  function hideLoadingAnimation() {

  }

  function addToTimeline(user, postContent, postID){
	const timeline = document.getElementById("timeline");
	const newPost = document.createElement("article");
	newPost.className = "flex flex-col shadow my-10 rounded-lg bg-white/25";
	newPost.innerHTML = `
	<div class="flex flex-col justify-start p-6">
		<p class="font-bold">@${user}</p>
		<p class="px-2">${postContent}</p>
    <p id=${postID}></p>
	</div>
	<div class="flex self-end gap-10 pb-6 mr-6">
		<i class="fa-regular fa-comment hover:text-cyan-600"></i>
		<i class="fa-solid fa-retweet hover:text-green-600"></i>
		<i class="fa-regular fa-heart hover:text-red-600"></i>
	</div>
  `;
  timeline.appendChild(newPost);
}

function useModerator(userPost, postID) {
// This is very uncertain with the default 0.9 threshold
const threshold = 0.7;
var postTag = document.getElementById('post-input');
showLoadingAnimation();
  userPost = postTag.value;
  console.log("Processing: " + userPost);

  //Receive response
  //Process based on policies
  // OK, post/mock
  // NOT OK, warning message
  // Control UI based on response

  const startTime = performance.now();

  toxicity.load(threshold)
    .then(model => {
      const sentences = userPost;
      model.classify(sentences)
        .then(predictions => {
          const endTime = performance.now();
          const executionTime = (endTime - startTime) / 1000;
          console.log(`Model execution time: ${executionTime.toFixed(2)} s`);
          handleModeratorReponse(predictions, postID);
        })
        .catch(error => {
          handleModeratorReponse("", postID);
          console.error("Error generating response", error);
        });
    })
    .catch(error => {
      handleModeratorReponse("", postID);
      console.error("Error loading model", error);
    });

};


function handleModeratorReponse(r, postID) {
  console.log(r);
  hideLoadingAnimation();
  const innerPost = document.getElementById(postID);
  innerPost.innerHTML = r.map(item => `<p>${JSON.stringify(item)}</p>`).join('');
}
