CREATE OR REPLACE FUNCTION insertlongurl(username character varying, password character varying, data json) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
       DECLARE
	nextId INTEGER;
	indexMapping TEXT[] := '{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","0","1","2","3","4","5","6","7","8","9"}';
	modval INTEGER;
	S TEXT;
	i TEXT;
	j JSON;
	country_block_length INTEGER;
	country_redirect_length INTEGER;
	device_type_length INTEGER;
	has_country_block BOOLEAN = false;
	has_device_redirect BOOLEAN = true;
       BEGIN
        loop
				SELECT  last_value + CASE WHEN is_called THEN 1 ELSE 0 END FROM urls_seq_seq INTO nextId;
				while nextId > 0 loop
	      	      		      modval := nextId % 62;
	      	 		      nextId := nextId / 62;
	      	    	 	      SELECT CONCAT(indexMapping[modval],S) INTO S;
			        end loop;
				BEGIN
					SELECT json_array_length(data->'country_block') INTO country_block_length;
					if country_block_length > 0 then
					   has_country_block = true;
					END if;
				
					if data->>'mobile_url' = '' AND data->>'desktop_url' = '' AND data->>'others_url' = '' then
					   has_device_redirect = false;
					END if;

					INSERT INTO urls(id,url,username,link_tag,password,not_found_url,android_deep_link,ios_deep_link,link_name,link_description,http_status,play_store_link,qr_code,country_block,device_type_redirect,mobile_url,desktop_url,expiration) VALUES(S,data->>'longUrl',username,data->>'tag',password,data->>'not_found_url',data->>'android_deep_link',data->>'ios_deep_link',data->>'name',data->>'description',CAST(data->>'http_status' AS INTEGER), data->>'play_store_link',CAST(data->>'qr_code' AS BOOLEAN),has_country_block,false, has_device_redirect,data->>'mobile_url',data->>'desktop_url',data->>'others_url',data->>'expiration');

					INSERT INTO campaign VALUES(S,data->'campaign'->>'name',data->'campaign'->>'medium',data->'campaign'->>'source',data->'campaign'->>'term',data->'campaign'->>'content',data->'campaign'->>'id');
					
					FOR i IN SELECT * from json_array_elements(data->'country_block')
					LOOP
						INSERT INTO country_block(id,country_code) VALUES(S,i);
					END LOOP;

					RETURN S;
					
				EXCEPTION WHEN unique_violation THEN
			  		  -- do nothing
				END;	
			
	END LOOP;
       COMMIT;
       END;
$$;
