INSERT INTO users(id, money)
	VALUES (p_user_id, p_funds)
	ON CONFLICT(id) DO
	    UPDATE SET money = users.money + p_funds
	RETURNING id,money;


