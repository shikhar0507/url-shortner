CREATE OR REPLACE FUNCTION insertLongUrl(url TEXT,username VARCHAR) RETURNS VARCHAR AS $$
       DECLARE
	nextId INTEGER;
	indexMapping TEXT[] := '{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","0","1","2","3","4","5","6","7","8","9"}';
	modval INTEGER;
	S TEXT;
	
       BEGIN
        loop
				
				SELECT  last_value + CASE WHEN is_called THEN 1 ELSE 0 END FROM urls_seq_seq INTO nextId;
				while nextId > 0 loop
	      	      		      modval := nextId % 62;
	      	 		      nextId := nextId / 62;
	      	    	 	      SELECT CONCAT(indexMapping[modval],S) INTO S;
			        end loop;
				BEGIN
					RAISE NOTICE '%',S;
					INSERT INTO urls(id,url,username) VALUES(S,url,username);
					RETURN S;
				EXCEPTION WHEN unique_violation THEN
			  		  -- do nothing
				END;
	
	END LOOP;
       COMMIT;
       END;
$$ LANGUAGE plpgsql;
