INSERT INTO features (name) VALUES ('football'), ('basketball'), ('tennis');
INSERT INTO tags (name) VALUES ('red'), ('blue'), ('green');
INSERT INTO banners (feature_id, data_id, is_active) VALUES 
(1, 101, true),
(2, 102, false),
(3, 103, true);
INSERT INTO banner_tags (banner_id, tag_id) VALUES (1, 1), (1, 2), (2, 1), (3, 3), (3, 2);
