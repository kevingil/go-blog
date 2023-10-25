	//Gather all resources
	//prompt will be reset after SSE connection close
	var prompt = "";

	// Htmx load
	htmx.onLoad(function (elt) {
		//stream request if completion page

		//parpse prompt if home page
		//stream if recipe page
		try {
			document.getElementById("generate-button").value = updatePrompt();
			document.getElementById('recipeform').addEventListener('change', function (event) {
				//home page
				document.getElementById("generate-button").value = updatePrompt();
				console.log("Updated -> " + updatePrompt());
			});
			//Bean Elevation slider
			document.getElementById('bean-elevation').addEventListener('input', function () {
				const roundedValue = Math.round(this.value / 100) * 100;
				document.getElementById('elevation-value').textContent = roundedValue + 'masl';
				this.value = roundedValue; // Snap the slider to the rounded value
			});

			console.log("Set prompt: " + getPrompt());

		} catch {
			let r = getPrompt();
			console.log("Requested: " + r);
			document.getElementById('result-placeholder')? requestRecipe(r) : console.log("no stream placeholder");
		}

	})


	//Swap listener
	document.addEventListener('htmx:afterSwap', function (event) {
		// Now you can use the 'data' variable as needed in your JavaScript code
		console.log("Swapped. Current prompt: ", getPrompt());
	});

	//Update and generate prompt
	//gathers info from the form and generates a prompt
	function updatePrompt() {
		try {
			let recipetype = document.querySelector('input[name="brewmethod"]:checked').value;
			let beanprocess = document.querySelector('input[name="beanprocess"]:checked').value;
			let beanelevtion = document.getElementById("elevation-value").innerText;
			let beancolor = document.querySelector('input[name="beancolor"]:checked').value;
			let beansvg = document.getElementById('beansvg');
			const parsedPrompt = `Drink type: ${recipetype}. `
				+ `Bean Process: ${beanprocess}. `
				+ `Growing elevation: ${beanelevtion}. `
				+ `Color: ${beancolor}. `;
			prompt = parsedPrompt;
			//change bean color
			if (beancolor === "Dark") {
				beansvg.style.fill = "#3d2b24";
			} else if (beancolor === "Medium-Dark") {
				beansvg.style.fill = "#553b32";
			} else if (beancolor === "Medium") {
				beansvg.style.fill = "#6d4d41";
			} else if (beancolor === "Medium-Light") {
				beansvg.style.fill = "#8c6354";
			} else if (beancolor === "Light") {
				beansvg.style.fill = "#ab7967";
			}
			return parsedPrompt;
		} catch (error) {
			console.log('No form available.');
			return 0;
		}
	}

	function getPrompt() {
		return prompt;
	}

	//copy prompt from UI

	//handle stream request
	//use prompt

	function requestRecipe(bean) {

		let msg = bean;
		let response = "";
		let sseConnection = null;

		if (sseConnection) {
			sseConnection.close();
		}

		sseConnection = new EventSource(`/api/stream-recipe?question=${encodeURIComponent(msg)}`);
		sseConnection.addEventListener("message", function (event) {
			document.getElementById('result-placeholder').style.display = 'none';
			response += event.data;
			renderMarkdown(response);

		});

		sseConnection.addEventListener("error", function (event) {
			sseConnection.close();
			sseConnection = null;
			console.log("Stream ended");
			showTryAgainButton();
		});
	}

	//prompt user to try again
	function showTryAgainButton() {
		const tryAgainElement = document.getElementById('tryagain');
		tryAgainElement.classList.remove('hidden');
	}

	// Convert stream to html as it arrives
	function renderMarkdown(markdown) {
		document.getElementById('result-placeholder').style.display = 'none';
		document.getElementById('result').innerHTML = marked.parse(markdown);
	}
