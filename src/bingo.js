document.getElementById("bingo").innerHTML = `
    <div class="container">
        <h1>BingoBingo</h1>
        <div class="form-group">
            <form id="bingo-form">
                <label for="name">Enter your name</label>
                <input type="text" id="name" name="name" class="form-control"><br>
                <label for="number">Enter a number between 1 and 10</label>
                <input type="text" id="number" name="number" class="form-control"><br>
                <input type="submit" value="Bingo!" class="btn btn-primary" /><br>
            </form>
        </div>
        <div id="result">
        </div>
    </div>
`;

const nameInput = document.getElementById("name");
let username = "";
nameInput.addEventListener("change", (event) => {
  username = event.target.value;
});

const numberInput = document.getElementById("number");
let number = "";
numberInput.addEventListener("change", (event) => {
  number = event.target.value;
});

document.getElementById("bingo-form")?.addEventListener("submit", (event) => {
  event.preventDefault();

  fetch("/api/try", {
    method: "POST",
    body: JSON.stringify({ "name": username, "number": number })
  })
    .then((res) => {
      if (res.ok) {
        return res.text()
          .then((data) => {
            const p = document.createElement("p");
            p.innerText = data;
            document.getElementById("result").append(p);
          })
        //const p = document.createElement("p");
        //p.innerText = res.text;
        //document.getElementById("bingo-form").append(p);
      }
    });
});
