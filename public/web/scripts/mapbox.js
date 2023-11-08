mapboxgl.accessToken = 'pk.eyJ1IjoiaGFyYW5zaXMiLCJhIjoiY2xhbXdmcWM1MGJnajN4bjA2OWQxZzV6ciJ9.ZjF8Laz8eHSYNvJOWNm93g';
const baseUrl = window.location.protocol + '//' + window.location.host + '/api/v1/assets/'
const map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/haransis/clesss239000801s0jzzy7nxz',
    center: [14.012033, 49.245659],
    zoom: 4,
});
map.on('load', () => {
    map.loadImage('icons/ic_marker.png',
    (error,image) => {
        if (error) throw error;

        map.addImage('geopoints-image', image);
        map.addSource('geopoints-source', {
            type: 'geojson',
            data: baseUrl + 'geojson.json'
        });
        map.addLayer({
            id: 'geopoints-layer',
            type: 'symbol',
            source: 'geopoints-source',
            'layout': {
                'icon-size': 0.35,
                'icon-image': 'geopoints-image',
                'text-field': ['get', 'name'],
                'text-font': [
                    'IBM Plex Mono Regular',
                    'IBM Plex Mono Bold'
                ],
                'text-size': 14,
                'text-offset': [1, 0],
                'text-anchor': 'left'
            }
        });

        map.on('click', 'geopoints-layer', (e) => {
            map.flyTo({
                center: e.features[0].geometry.coordinates
            });
        });
        map.on('mouseenter', 'geopoints-layer', () => {
            map.getCanvas().style.cursor = 'pointer';
        });
        map.on('mouseleave', 'geopoints-layer', () => {
            map.getCanvas().style.cursor = '';
        });
    })
});