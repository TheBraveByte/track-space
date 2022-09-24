//to fecth data from json file
let data;
fetch('../data.json')
    .then(res => res.json())
    .then(jsondata => data = jsondata)
    .then(() => {
        // set the dimensions and margins of the graph
        let width = 500;
        let height = 500;
        let margin = 40;

        // The radius of the pie-plot is half the width or half the height (the smallest one). I subtract a bit of margin.
        let radius = Math.min(width, height) / 2 - margin

        // append the svg object to the div called 'doughnut'
        let svg = d3.select("#doughnut")
            .append("svg")
            .attr("width", width)
            .attr("height", height)
            .append("g")
            .attr("transform", "translate(" + width / 2 + "," + height / 2 + ")");

        // Create dummy data
        //let data = {a: 9, b: 20, c:30, d:8, e:12}

        // set the color scale
        let color = d3.scaleOrdinal()
            .domain(data)
            .range(d3.schemeSet2);

        // Compute the position of each group on the pie:
        let pie = d3.pie()
            .value(function(d) {return d["total"]; })
        let data_ready = pie(data)
        // data_ready is an array of 17 objects

        // shape helper to build arcs:
        let arcGenerator = d3.arc()
            .innerRadius(100)
            .outerRadius(radius)
            .cornerRadius(10);

        // Build the pie chart: Basically, each part of the pie is a path that we build using the arc function.
        svg.selectAll('mySlices')
            .data(data_ready)
            .enter()
            .append('path')
            .attr('d', arcGenerator)
            .attr('fill', function(d){ return(color(d.data.total)) })
            .attr("stroke", "white")
            .style("stroke-width", "2px")
            .style("opacity", 0.9)
            .append("title")
            .text((item) => `${item.data.total} total projects on ${item.data.date}`);


        // Now add the annotation. Use the centroid method to get the best coordinates
        svg.selectAll('mySlices')
            .data(data_ready)
            .enter()
            .append('text')
            .text(function(d){ return `${d.data.total} total`})
            .attr("transform", function(d) { return "translate(" + arcGenerator.centroid(d) + ")";  })
            .style("text-anchor", "middle")
            .style('font-family', 'Overpass')
            .style("font-size", 10)
            .style('font-weight', 900)
            .attr('fill', "white")
    })
