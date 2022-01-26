library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use ieee.numeric_std.all;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}; output_size: integer:={{.OutputSize}}); 
Port (
clk : in std_logic;
rst : in std_logic;
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (output_size-1 downto 0));
end {{.EntityName}};

architecture Behavioral of {{.EntityName}} is
    signal acc : unsigned(output_size-1 downto 0) := (others => '0');
begin
    
    p_asynchronous_reset : process(clk, rst) is
    begin
        if rst = '1' then                  
            acc <= (others => '0');
        elsif rising_edge(clk) then       
            acc <= acc + (unsigned(a)*unsigned(b));
        end if;
    end process p_asynchronous_reset;

    prod <= STD_LOGIC_VECTOR(acc);
end Behavioral;