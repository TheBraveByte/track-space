import data from "./data.js";
//to set the width, height and padding of the svg
const width = 500;
const height = 300;
const padding = 40;

//to set the x scale/axis

const dateParser = d3.timeParse("%Y-%m-%d"),
    formatDate = d3.timeFormat("%b %d"),
    formatMonth = d3.timeFormat("%b");
const xAccessor = (d) => dateParser(d["date"]);

const xScale = d3.scaleTime()
    .domain(d3.extent(data, xAccessor))
    .range([padding, width - 50]);
const xAxis = d3.axisBottom(xScale);

//to set the y scales/axis
const yScale = d3.scaleLinear()
    .domain([0, d3.max(data, (item) => item["total"])])
    .range([height - padding, padding]);
const yAxis = d3.axisLeft(yScale);

//appending svg tag to the chart div
const svg = d3.select(".chart")
    .append("svg")
    .attr("width", width)
    .attr("height", height)
    .style("background-color", "#ecfdf5")
    .style("border-radius", "10px")
    .style("border", "1px solid lightgrey")

//setting the x and y-axis
svg.append("g")
    .attr("transform", "translate(0," + (height - padding) + ")")
    .call(xAxis);
svg.append("g")
    .attr("transform", "translate(" + padding + ", 0)")
    .call(yAxis);

svg.selectAll("circle")
    .data(data)
    .enter()
    .append("circle")
    .attr("cx", (item) => xScale(xAccessor(item)))
    .attr("cy", (item) => yScale(item["total"]))
    .attr("r", (item) => 3.5)
    .attr("class", "point")
    .append("title")
    .text((item) => `Article: ${item["article"]}, Code: ${item["code"]}, Text: ${item["text"]}`);

// line chart title
svg.append('text')
    .attr('x', width / 2)
    .attr('y', 25)
    .attr('text-anchor', 'middle')
    .style('font-family', 'Overpass')
    .style('font-size', 16)
    .text('Line Chart');


let line = d3.line()
    .x((item) => xScale(xAccessor(item)))
    .y((item) => yScale(item["total"]))
//.curve(d3.curveMonotoneX)

svg.append("path")
    .datum(data)
    .attr("d", line)
    .style("fill", "none")
    .style("stroke", "#059669")
    .style("stroke-width", "2")
