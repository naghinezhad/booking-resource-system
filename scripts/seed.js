const RESOURCES = 100;
const RESERVATIONS = 10000;
const BATCH_SIZE = 1000;

const resources = [];

for (let i = 1; i <= RESOURCES; i++) {
  resources.push({ name: `resource-${i}` });
}

const result = db.resources.insertMany(resources);
const resourceIds = Object.values(result.insertedIds);

print("Inserted resources:", resourceIds.length);

function randomInt(max) {
  return Math.floor(Math.random() * max);
}

let batch = [];
let inserted = 0;

for (let i = 0; i < RESERVATIONS; i++) {
  const resourceId = resourceIds[randomInt(resourceIds.length)];

  const startOffset = randomInt(7 * 24 * 60); // up to 7 days
  const duration = (randomInt(4) + 1) * 30; // 30-120 min

  const start = new Date(Date.now() + startOffset * 60000);
  const end = new Date(start.getTime() + duration * 60000);

  batch.push({
    resource_id: resourceId,
    start_time: start,
    end_time: end,
    created_at: new Date(),
  });

  if (batch.length === BATCH_SIZE) {
    db.reservations.insertMany(batch);
    inserted += batch.length;
    batch = [];
  }
}

if (batch.length) {
  db.reservations.insertMany(batch);
  inserted += batch.length;
}

print("Inserted reservations:", inserted);
