let datum;

fetch('/static/json/stat.json')
    .then(res => res.json())
    .then(jsondata => datum = jsondata)
    .then(() => {
        document.getElementById("total-users").innerHTML = datum.length;

        //to get the sum total of values by keys
        const getKey = (arr, key) => {
            return arr.reduce((accumulator, current) => accumulator + Number(current[key]), 0)
        }

        //total total-projects
        const total = getKey(datum, "total")
        document.getElementById("total-projects").innerHTML = total;

        //total todo
        const todo = getKey(datum, "todo")
        document.getElementById("total-todo").innerHTML = todo;
    })