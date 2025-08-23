create function trim_crypto_history()
returns trigger as $$
begin

    DELETE FROM prices
    WHERE id NOT IN (
        SELECT id FROM prices
        WHERE crypto = NEW.crypto
        ORDER BY updated_at ASC
        LIMIT 100
    );

   return new;
end;
$$ language plpgsql;


create trigger crypto_history_limit
after insert on prices
for each row
execute function trim_crypto_history();


