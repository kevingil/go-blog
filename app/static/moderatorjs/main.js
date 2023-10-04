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

document.getElementById('post-input').addEventListener('keypress', function(event) {
  if (event.keyCode === 13) { 
      event.preventDefault(); 
      document.getElementById('postButton').click(); 
  }
});


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
  
  function hideLoadingAnimation(postID) {
    var spinnerID = postID + "_spinner";
    var spinnerElement = document.getElementById(spinnerID);
    spinnerElement.style.display = "none";
  }

  function addToTimeline(user, postContent, postID){
	const timeline = document.getElementById("timeline");
	const newPost = document.createElement("article");
	newPost.className = "flex flex-col shadow my-10 rounded-lg bg-white/25";
	newPost.innerHTML = `
	<div class="flex flex-col justify-start p-6">
		<p class="font-bold">@${user}</p>
		<p class="px-2">${postContent}</p>
    <p id="${postID}" class="mt-4"></p>
    <div class="text-center">
    <div id="${postID}_spinner" role="status">
        <svg aria-hidden="true" class="inline w-8 h-8 mr-2 text-gray-200 animate-spin dark:text-gray-400 fill-blue-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor"/>
            <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill"/>
        </svg>
        <span class="">Running model...</span>
    </div>
</div>
	</div>
	<div class="flex self-end gap-10 pb-6 mr-6">
		<i class="fa-regular fa-comment hover:text-cyan-600"></i>
		<i class="fa-solid fa-retweet hover:text-green-600"></i>
		<i class="fa-regular fa-heart hover:text-red-600"></i>
	</div>
  `;
  timeline.insertBefore(newPost, timeline.firstChild);
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
          handleModeratorResponse(predictions, postID, startTime);
        })
        .catch(error => {
          handleModeratorResponse("", postID, startTime);
          console.error("Error generating response", error);
        });
    })
    .catch(error => {
      handleModeratorResponse("", postID, startTime);
      console.error("Error loading model", error);
    });

};


function handleModeratorResponse(response, postID, startTime) {
  const endTime = performance.now();
  const executionTime = ((endTime - startTime) / 1000).toFixed(2);
  hideLoadingAnimation(postID);
  const innerPost = document.getElementById(postID);

  //Handle error
  if (!Array.isArray(response)) {
    innerPost.innerHTML = "<p>Invalid reponse, something's wrong.</p>";
    return;
  }

  //Parse reponse
  const html = response.map(item => {
    let prob0 = ((item.results[0].probabilities[0])*100).toFixed(2);
    let prob1 = ((item.results[0].probabilities[1])*100).toFixed(2);
    return `
      <p class="uppercase font-semibold">${item.label}</p>
      <p class="pb-4">${item.results[0].match ? `<span class="text-red-700">YES</span>, confidence: ${prob1 + "%"}` : `<span class="text-green-700">NO</span>, confidence ${prob0 + "%"}`}</p>
    `;
  }).join('');

  innerPost.innerHTML = html + `<p class="py-4 text-lg">Model Finished in ${ executionTime }s</p>`;
  
}
