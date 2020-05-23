(SELECT f.species_name,
        AVG(f.height) AS average_height, AVG(f.diameter) AS average_diameter
 FROM flora AS f
 WHERE f.species_name = 'Banksia'
    OR f.species_name = 'Sheoak'
    OR f.species_name = 'Wattle'
 GROUP BY f.species_name, f.observation_date)

UNION ALL

(SELECT b.species_name,
        AVG(b.height) AS average_height, AVG(b.diameter) AS average_diameter
 FROM botanic_garden_flora AS b
 WHERE b.species_name = 'Banksia'
    OR b.species_name = 'Sheoak'
    OR b.species_name = 'Wattle'
 GROUP BY b.species_name, b.observation_date);