library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use ieee.numeric_std.all;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}); -- Keep at 2
Port ( 
a : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
b : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};

architecture Behavioral of {{.EntityName}} is
begin
    prod <= STD_LOGIC_VECTOR(unsigned(a) * unsigned(b));
end Behavioral;