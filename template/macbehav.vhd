library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use IEEE.NUMERIC_STD.ALL;

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
    signal mult_output : std_logic_vector(word_size * 2 - 1 downto 0) := (others => '0');
begin

mult: entity work.{{.Multiplier.EntityName}} port map(A=>A, B=>B,prod=>mult_output);


p_asynchronous_reset : process(clk, rst) is
begin
   if rst = '1' then                  
      acc <= (others => '0');
   elsif rising_edge(clk) then       
      acc <= acc + unsigned(mult_output);
   end if;
end process p_asynchronous_reset;

prod <= std_logic_vector(acc);

end Behavioral;