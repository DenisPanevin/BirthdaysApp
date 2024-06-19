CREATE TABLE if not exists subscriptions (
                                             subscription_id SERIAL PRIMARY KEY,
                                             subscriber_id INT,
                                             subscribed_to_id INT,
                                             FOREIGN KEY (subscriber_id) REFERENCES users(id),
                                             FOREIGN KEY (subscribed_to_id) REFERENCES users(id),
                                             UNIQUE(subscriber_id, subscribed_to_id)
);