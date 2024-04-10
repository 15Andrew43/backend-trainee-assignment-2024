SELECT data_id, feature_id, tag_id, is_active, created_at, updated_at
   FROM banners b
   INNER JOIN banner_tags bt ON b.id = bt.banner_id
;

-- SELECT * from banners;