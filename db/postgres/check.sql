SELECT *
   FROM banners b
   INNER JOIN banner_tags bt ON b.id = bt.banner_id
;

-- SELECT * from banners;