import * as d3 from 'd3';

async function fetchDBInfo() {
  const resp = await fetch('/api/db/info', {
    method: 'POST',
  });

  return resp.json();
}

function convertDbInfoData(data) {
  const firstScheme = data.schemes[0]; // todo: handle other schemes

  const groups = firstScheme.tables.map(({name, columns}) => ({
    x: 50, // todo compute coordinates?
    y: 100,
    name,
    columns,
  }));
  const links = firstScheme.relationships.map((rel) => ({from: rel.from, to: rel.to}));
  return [groups, links];
}

async function app() {
  const info = await fetchDBInfo();
  const [groups, links] = convertDbInfoData(info);

  const rectHeight = 20;
  const rectWidth = 120;
  const radius = 5;

  const svg = d3.select('body').append('svg')
    .attr('width', window.innerWidth - 10)
    .attr('height', window.innerHeight - 10)
    .style('background-color', '#603db0')
    .attr('stroke-width', 1);

  svg.append('path')
    .attr('stroke', 'red')
    .attr('fill', 'none');

  svg.selectAll('rect')
    .data(groups)
    .join('g')
    .call(d3.drag().on('drag', dragged));

  groups.forEach((data, index) => {
    const node = svg.selectAll('g').filter((d, i) => i === index);

    node.append('rect')
      .attr('x', (d) => d.x)
      .attr('y', (d) => d.y)
      .attr('width', rectWidth)
      .attr('height', rectHeight)
      .attr('fill', '#a8c5f5')
      .attr('stroke', '#000000');

    node.append('text')
      .attr('x', (d) => d.x + ((rectWidth / 2) - 5.5 * d.name.length / 2))
      .attr('y', (d) => d.y + 13)
      .text((d) => d.name)
      .style('font', "9px 'Comic Sans MS'")
      .style('fill', 'black'); // заливка текста цветом

    data.columns.forEach((column, columnI) => {
      node.append('rect')
        .attr('x', (d) => d.x)
        .attr('y', (d) => d.y + (columnI + 1) * rectHeight)
        .attr('width', rectWidth)
        .attr('height', rectHeight)
        .attr('fill', '#ffffff')
        .attr('stroke', '#000000');

      node.append('text')
        .attr('x', (d) => d.x + 10)
        .attr('y', (d) => d.y + 13 + (columnI + 1) * rectHeight)
        .text(`${column.name} [${column.type.toUpperCase()}]`)
        .style('font', "9px 'Comic Sans MS'")
        .style('fill', 'black');

      node.append('circle')
        .attr('cx', (d) => d.x + rectWidth - 10)
        .attr('cy', (d) => d.y + (columnI + 1) * rectHeight + 10)
        .attr('r', radius)
        .attr('fill', 'white')
        .attr('stroke', 'green')
        .on('click', addLink);
    });
  });

  svg.selectAll('path')
    .data(links).join('path')
    .attr('stroke', 'black')
    .attr('fill', 'none')
    .attr('d', (d) => getPath(d));

  function getPath(link = {}) {
    const fromTable = groups.find((item) => item.name === link.from.table);
    const toTable = groups.find((item) => item.name === link.to.table);

    const indexColumnFromTable = fromTable.columns.findIndex((item) => item.name === link.from.column);
    const indexColumnToTable = toTable.columns.findIndex((item) => item.name === link.to.column);

    let fromX = fromTable.x;
    const fromY = fromTable.y + (indexColumnFromTable + 1) * 20 + 10;

    let toX = toTable.x;
    let toY = toTable.y + (indexColumnToTable + 1) * 20 + 10;

    const fixToX = toX + rectWidth < fromX ? rectWidth : 0;
    const fixFromX = fromX + rectWidth < toX ? rectWidth : 0;

    fromX += fixFromX;
    toX += fixToX;

    const Q = {
      x: (fromX + (toX - fromX) / 2) - 20,
      y: fromY + (toY - fromY),
    };

    if (link.to.point) {
      toX = link.to.point.x;
      toY = link.to.point.y;
    }

    return `M${fromX} ${fromY} Q ${Q.x} ${Q.y} ${toX} ${toY}`;
  }

  function addLink(e, d) {
    if (svg.on('mousemove')) {
      console.log('off');
      const index = Math.round((e.pageY - d.y) / rectHeight) - 2;
      console.log('table', d.name);
      console.log('column', d.columns[index].name);

      if (links[links.length - 1].from.table === d.name) {
        links.pop();
      } else {
        links[links.length - 1].to = {
          table: d.name,
          column: d.columns[index].name,
        };
      }

      svg.selectAll('path')
        .data(links).join('path')
        .attr('stroke', 'black')
        .attr('fill', 'none')
        .attr('d', (d) => getPath(d));
      svg.on('mousemove', undefined);
    } else {
      console.log('on');
      const index = Math.round((e.pageY - d.y) / rectHeight) - 2;

      links.push({
        from: {table: d.name, column: d.columns[index].name},
        to: {
          table: 'test',
          column: 'id',
          point: {x: e.pageX, y: e.pageY},
        },
      });

      svg.selectAll('path')
        .data(links).join('path')
        .attr('stroke', 'black')
        .attr('fill', 'none')
        .attr('d', (d) => getPath(d));
      svg.on('mousemove', eventAdd);
    }
  }

  function eventAdd(e, d) {
    links[links.length - 1].to.point.x = e.pageX;
    links[links.length - 1].to.point.y = e.pageY;
    svg.selectAll('path')
      .data(links).join('path')
      .attr('stroke', 'black')
      .attr('fill', 'none')
      .attr('d', (d) => getPath(d));
  }

  function dragged(event) {
    d3.select(this).select('rect').attr('x', (d) => {
      d.x += event.dx;
      return d.x;
    }).attr('y', (d) => {
      d.y += event.dy;
      return d.y;
    });

    d3.select(this).selectAll('rect').filter((d, i) => i !== 0)
      .attr('x', (d) => d.x)
      .attr('y', (d, i) => d.y + 20 * (i + 1));

    d3.select(this).select('text').attr('x', (d) => d.x + ((rectWidth / 2) - 5.5 * d.name.length / 2)).attr('y', (d, i) => d.y + 13 + (i * 20));

    d3.select(this).selectAll('text').filter((d, i) => i !== 0)
      .attr('x', (d) => d.x + 10)
      .attr('y', (d, i) => d.y + 13 + ((i + 1) * 20));

    d3.select(this).selectAll('circle')
      .attr('cx', (d) => d.x + rectWidth - 10)
      .attr('cy', (d, i) => d.y + (i + 1) * rectHeight + 10);

    svg.selectAll('path').attr('d', (d) => getPath(d));
  }
}

app().then();
