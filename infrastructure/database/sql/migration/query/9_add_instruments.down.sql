-- Down Migration
DELETE FROM instrument WHERE name IN (
  'Harp',
  'Lute',
  'Fiddle',
  'Violin',
  'Viola',
  'Cello',
  'Double Bass',
  'Electric Guitar',
  'Flute',
  'Oboe',
  'Clarinet',
  'Fife',
  'Panpipes',
  'Saxophone',
  'Trumpet',
  'Trombone',
  'Tuba',
  'Horn',
  'Piano',
  'Timpani',
  'Bongo',
  'Bass Drum',
  'Snare Drum',
  'Cymbal'
);
